package v1

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

func TestInsertMany_AllNewRecords(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Insert 3 new records
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1",
			CreatedAt:   timestamppb.New(baseTime),
			UpdatedAt:   timestamppb.New(baseTime),
			DisplayName: "Dog",
			Description: "A friendly dog",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    9.5,
			Price:       500.0,
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(baseTime),
			UpdatedAt:   timestamppb.New(baseTime),
			DisplayName: "Cat",
			Description: "A curious cat",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    9.0,
			Price:       400.0,
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3",
			CreatedAt:   timestamppb.New(baseTime),
			UpdatedAt:   timestamppb.New(baseTime),
			DisplayName: "Bird",
			Description: "A singing bird",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    8.5,
			Price:       300.0,
		}.Build(),
	}

	query, params, err := pgdb_v1.InsertMany(pgdb_v1.DialectV13, pets...)
	require.NoError(t, err)

	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	require.Equal(t, int64(3), res.RowsAffected())
}

func TestInsertMany_OneConflictAcceptNew(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Insert initial record
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	initialPet := Pet_builder{
		TenantId:    "t1",
		Id:          "pet1",
		CreatedAt:   timestamppb.New(oldTime),
		UpdatedAt:   timestamppb.New(oldTime),
		DisplayName: "Old Name",
		Description: "Old description",
		Elapsed:     durationpb.New(time.Hour),
		Profile:     &structpb.Struct{},
		Cuteness:    5.0,
		Price:       100.0,
	}.Build()

	query, params, err := pgdb_v1.Insert(initialPet, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	// Insert with one conflicting record (newer timestamp) and two new records
	newTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	pets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1", // Conflict - should update
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New Name",
			Description: "New description",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    9.5,
			Price:       500.0,
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2", // New
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "Cat",
			Description: "A curious cat",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    9.0,
			Price:       400.0,
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3", // New
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "Bird",
			Description: "A singing bird",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    8.5,
			Price:       300.0,
		}.Build(),
	}

	query, params, err = pgdb_v1.InsertMany(pgdb_v1.DialectV13, pets...)
	require.NoError(t, err)

	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	// RowsAffected counts inserted + updated rows (1 updated, 2 inserted = 3)
	require.Equal(t, int64(3), res.RowsAffected())

	// Verify the first record was updated
	var displayName string
	var updatedAt string
	err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name", "pb$updated_at" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", "pet1").Scan(&displayName, &updatedAt)
	require.NoError(t, err)
	require.Equal(t, "New Name", displayName)
}

func TestInsertMany_OneConflictRejectNew(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Insert initial record with newer timestamp
	newTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	initialPet := Pet_builder{
		TenantId:    "t1",
		Id:          "pet1",
		CreatedAt:   timestamppb.New(newTime),
		UpdatedAt:   timestamppb.New(newTime),
		DisplayName: "Newer Name",
		Description: "Newer description",
		Elapsed:     durationpb.New(time.Hour),
		Profile:     &structpb.Struct{},
		Cuteness:    9.5,
		Price:       500.0,
	}.Build()

	query, params, err := pgdb_v1.Insert(initialPet, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	// Try to insert with one conflicting record (older timestamp) and two new records
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1", // Conflict - should NOT update (older)
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old Name",
			Description: "Old description",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    5.0,
			Price:       100.0,
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2", // New
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Cat",
			Description: "A curious cat",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    9.0,
			Price:       400.0,
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3", // New
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Bird",
			Description: "A singing bird",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
			Cuteness:    8.5,
			Price:       300.0,
		}.Build(),
	}

	query, params, err = pgdb_v1.InsertMany(pgdb_v1.DialectV13, pets...)
	require.NoError(t, err)

	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	// RowsAffected: 2 new inserts, 1 conflict that was NOT updated (WHERE clause failed)
	// PostgreSQL counts inserts + successful updates, but not updates blocked by WHERE
	require.Equal(t, int64(2), res.RowsAffected())

	// Verify the first record was NOT updated (kept the newer one)
	var displayName string
	err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", "pet1").Scan(&displayName)
	require.NoError(t, err)
	require.Equal(t, "Newer Name", displayName, "Should keep the existing newer record")
}

func TestInsertMany_MultipleConflictsMixed(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Insert initial records with different timestamps
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)

	initialPets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1",
			CreatedAt:   timestamppb.New(t1),
			UpdatedAt:   timestamppb.New(t1),
			DisplayName: "Dog v1",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(t3),
			UpdatedAt:   timestamppb.New(t3),
			DisplayName: "Cat v3",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	query, params, err := pgdb_v1.InsertMany(pgdb_v1.DialectV13, initialPets...)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	// Insert batch with:
	// - pet1: newer (should update)
	// - pet2: older (should NOT update)
	// - pet3: new (should insert)
	pets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1", // Conflict - newer, should update
			CreatedAt:   timestamppb.New(t2),
			UpdatedAt:   timestamppb.New(t2),
			DisplayName: "Dog v2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2", // Conflict - older, should NOT update
			CreatedAt:   timestamppb.New(t2),
			UpdatedAt:   timestamppb.New(t2),
			DisplayName: "Cat v2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3", // New
			CreatedAt:   timestamppb.New(t2),
			UpdatedAt:   timestamppb.New(t2),
			DisplayName: "Bird v2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	query, params, err = pgdb_v1.InsertMany(pgdb_v1.DialectV13, pets...)
	require.NoError(t, err)

	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	// RowsAffected: 1 update (pet1), 1 insert (pet3), 1 rejected update (pet2)
	require.Equal(t, int64(2), res.RowsAffected())

	// Verify results
	var displayName string

	// pet1 should be updated to v2
	err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", "pet1").Scan(&displayName)
	require.NoError(t, err)
	require.Equal(t, "Dog v2", displayName, "pet1 should be updated to v2")

	// pet2 should still be v3 (not updated to v2)
	err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", "pet2").Scan(&displayName)
	require.NoError(t, err)
	require.Equal(t, "Cat v3", displayName, "pet2 should remain v3 (newer)")

	// pet3 should be inserted
	err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", "pet3").Scan(&displayName)
	require.NoError(t, err)
	require.Equal(t, "Bird v2", displayName, "pet3 should be inserted")
}

func TestInsertMany_AllConflictsAcceptAll(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Insert initial records
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	initialPets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1",
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old 1",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old 2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3",
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old 3",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	query, params, err := pgdb_v1.InsertMany(pgdb_v1.DialectV13, initialPets...)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	// Update all records with newer timestamp
	newTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	updatedPets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1",
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New 1",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New 2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3",
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New 3",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	query, params, err = pgdb_v1.InsertMany(pgdb_v1.DialectV13, updatedPets...)
	require.NoError(t, err)

	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	require.Equal(t, int64(3), res.RowsAffected())

	// Verify all records were updated
	for i := 1; i <= 3; i++ {
		var displayName string
		petID := "pet" + strconv.Itoa(i)
		expectedName := "New " + strconv.Itoa(i)
		err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", petID).Scan(&displayName)
		require.NoError(t, err)
		require.Equal(t, expectedName, displayName)
	}
}

func TestInsertMany_AllConflictsRejectAll(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Insert initial records with newer timestamp
	newTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	initialPets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1",
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New 1",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New 2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3",
			CreatedAt:   timestamppb.New(newTime),
			UpdatedAt:   timestamppb.New(newTime),
			DisplayName: "New 3",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	query, params, err := pgdb_v1.InsertMany(pgdb_v1.DialectV13, initialPets...)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	// Try to update all records with older timestamp (should be rejected)
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedPets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet1",
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old 1",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old 2",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet3",
			CreatedAt:   timestamppb.New(oldTime),
			UpdatedAt:   timestamppb.New(oldTime),
			DisplayName: "Old 3",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	query, params, err = pgdb_v1.InsertMany(pgdb_v1.DialectV13, updatedPets...)
	require.NoError(t, err)

	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	// All 3 updates were rejected by WHERE clause, so 0 rows affected
	require.Equal(t, int64(0), res.RowsAffected())

	// Verify no records were updated (kept the newer ones)
	for i := 1; i <= 3; i++ {
		var displayName string
		petID := "pet" + strconv.Itoa(i)
		expectedName := "New " + strconv.Itoa(i)
		err = pg.DB.QueryRow(ctx, `SELECT "pb$display_name" FROM pb_pet_models_animals_v1_8a3723d5 WHERE "pb$tenant_id" = $1 AND "pb$id" = $2`, "t1", petID).Scan(&displayName)
		require.NoError(t, err)
		require.Equal(t, expectedName, displayName, "Should keep the newer record")
	}
}
