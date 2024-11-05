package mongodb

import (
	"context"
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dateKey = "createdAt"

type Repository[DType any, SType any] struct {
	Collection    *mongo.Collection
	FilterKeys    []string
	AggregatePipe []bson.M
}

func (repo *Repository[DType, SType]) GetById(ctx context.Context, id primitive.ObjectID) (*DType, *gerror.Error) {
	//obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": id,
	}
	return repo.GetByQuery(ctx, query)
}

func (repo *Repository[DType, SType]) GetByDna(ctx context.Context, id primitive.ObjectID, dna string) (*DType, *gerror.Error) {
	if len(dna) == 0 {
		return nil, &gerror.Error{
			Code: "api.user.nf",
		}
	}
	query := bson.M{
		"_id": id,
		"dna": primitive.Regex{Pattern: "^" + dna, Options: "i"},
	}
	return repo.GetByQuery(ctx, query)
}

func (repo *Repository[DType, SType]) GetNullable(ctx context.Context, id *primitive.ObjectID) (*DType, *gerror.Error) {
	if id == nil {
		return nil, &gerror.Error{
			Code: "api.user.nf",
		}
	}
	query := bson.M{
		"_id": id,
	}
	return repo.GetByQuery(ctx, query)
}

func (repo *Repository[DType, SType]) GetByQuery(ctx context.Context, query bson.M) (*DType, *gerror.Error) {
	var item DType
	res := repo.Collection.FindOne(ctx, query)
	if err := res.Decode(&item); err != nil {
		return nil, gerror.GetError(err)
	}
	return &item, nil
}

func (repo *Repository[DType, SType]) GetDetail(ctx context.Context, id primitive.ObjectID) (*DType, *gerror.Error) {
	return repo.GetDetailByQuery(ctx, bson.M{"_id": id})
}
func (repo *Repository[DType, SType]) GetDetailByQuery(ctx context.Context, query bson.M) (*DType, *gerror.Error) {

	pipeline := append(repo.AggregatePipe, bson.M{"$match": query})
	pipeline = append(pipeline, bson.M{"$limit": 1})
	pipeline = append(pipeline, bson.M{"$skip": 0})
	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		var item DType
		err = cur.Decode(&item)
		if err != nil {
			return nil, gerror.GetError(err)
		}
		return &item, nil
	}
	return nil, &gerror.Error{
		Code:   "mongo.notfound",
		Detail: "Document not found",
	}
}

func (repo *Repository[DType, SType]) GetByFirst(ctx context.Context, key string, value any) (*DType, *gerror.Error) {
	query := bson.M{
		key: value,
	}
	return repo.GetByQuery(ctx, query)
}

func (repo *Repository[DType, SType]) List(ctx context.Context, filters response.ListRequest) (*response.ListResponse[DType], *gerror.Error) {
	query := bson.M{}
	gte, gteOk := filters.Filters["gte"]
	lte, lteOk := filters.Filters["lte"]
	if gteOk && lteOk {
		query[dateKey] = bson.M{"$gte": gte, "$lt": lte}
	} else if lteOk {
		query[dateKey] = bson.M{"$lt": lte}
	} else if gteOk {
		query[dateKey] = bson.M{"$gte": gte}
	}

	for key, val := range filters.Filters {
		if key == "gte" || key == "lte" {
			continue
		}
		if key == "keyword" {
			if len(repo.FilterKeys) == 0 {
				continue
			}
			and := bson.A{}
			for _, filterKey := range repo.FilterKeys {
				and = append(and, bson.M{filterKey: primitive.Regex{Pattern: val.(string), Options: "i"}})
			}
			query["$or"] = and
		} else if key == "dna" {
			query[key] = primitive.Regex{Pattern: "^" + val.(string), Options: "i"}
		} else {
			query[key] = val
		}
	}
	pipeline := append(repo.AggregatePipe, bson.M{"$match": query})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		dateKey: -1,
	}})
	pipeline = append(pipeline, bson.M{"$skip": int64((filters.Page - 1) * filters.PageSize)})
	pipeline = append(pipeline, bson.M{"$limit": filters.PageSize})

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	defer cur.Close(ctx)
	var list []DType
	for cur.Next(ctx) {
		var item DType
		err := cur.Decode(&item)
		if err != nil {
			return nil, gerror.GetError(err)
		}
		list = append(list, item)
	}
	if err := cur.Err(); err != nil {
		return nil, gerror.GetError(err)
	}
	count, err := repo.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	if len(list) == 0 {
		list = []DType{}
	}
	return &response.ListResponse[DType]{
		Page:     filters.Page,
		PageSize: filters.PageSize,
		Total:    count,
		List:     list,
	}, nil
}

func (repo *Repository[DType, SType]) ListBasic(ctx context.Context, filters response.ListRequest) (*response.ListResponse[SType], *gerror.Error) {
	query := bson.M{}
	gte, gteOk := filters.Filters["gte"]
	lte, lteOk := filters.Filters["lte"]
	if gteOk && lteOk {
		query[dateKey] = bson.M{"$gte": gte, "$lt": lte}
	} else if lteOk {
		query[dateKey] = bson.M{"$lt": lte}
	} else if gteOk {
		query[dateKey] = bson.M{"$gte": gte}
	}

	for key, val := range filters.Filters {
		if key == "gte" || key == "lte" {
			continue
		}
		if key == "keyword" {
			if len(repo.FilterKeys) == 0 {
				continue
			}
			and := bson.A{}
			for _, filterKey := range repo.FilterKeys {
				and = append(and, bson.M{filterKey: primitive.Regex{Pattern: val.(string), Options: "i"}})
			}
			query["$or"] = and
		} else if key == "dna" {
			query[key] = primitive.Regex{Pattern: "^" + val.(string), Options: "i"}
		} else {
			query[key] = val
		}
	}
	//skip := int64((filters.Page - 1) * filters.PageSize)
	var pipeline []bson.M
	pipeline = append(pipeline, bson.M{"$match": query})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		dateKey: -1,
	}})
	pipeline = append(pipeline, bson.M{"$skip": int64((filters.Page - 1) * filters.PageSize)})
	pipeline = append(pipeline, bson.M{"$limit": filters.PageSize})

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	defer cur.Close(ctx)
	var list []SType
	for cur.Next(ctx) {
		var item SType
		err := cur.Decode(&item)
		if err != nil {
			return nil, gerror.GetError(err)
		}
		list = append(list, item)
	}
	if err := cur.Err(); err != nil {
		return nil, gerror.GetError(err)
	}
	count, err := repo.Collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, gerror.GetError(err)
	}

	if len(list) == 0 {
		list = []SType{}
	}
	return &response.ListResponse[SType]{
		Page:     filters.Page,
		PageSize: filters.PageSize,
		Total:    count,
		List:     list,
	}, nil
}

func (repo *Repository[DType, SType]) Add(ctx context.Context, schema DType) (*DType, *gerror.Error) {
	res, err := repo.Collection.InsertOne(ctx, schema)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	id := res.InsertedID.(primitive.ObjectID)
	return repo.GetById(ctx, id)
}

func (repo *Repository[DType, SType]) Update(ctx context.Context, id primitive.ObjectID, schema interface{}) (bool, *gerror.Error) {
	opts := options.Update().SetUpsert(false)
	_, err := repo.Collection.UpdateOne(
		ctx,
		bson.D{{"_id", id}},
		bson.D{{"$set", schema}},
		opts)
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}

func (repo *Repository[DType, SType]) Increment(ctx context.Context, id primitive.ObjectID, key string, val int) (bool, *gerror.Error) {
	opts := options.Update().SetUpsert(false)
	_, err := repo.Collection.UpdateOne(
		ctx,
		bson.D{{"_id", id}},
		bson.D{{"$inc", bson.M{
			key: val,
		}}},
		opts)
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}

func (repo *Repository[DType, SType]) Replace(ctx context.Context, id primitive.ObjectID, schema DType) (*DType, *gerror.Error) {
	opts := options.Replace().SetUpsert(false)
	_, err := repo.Collection.ReplaceOne(
		ctx,
		bson.D{{"_id", id}},
		bson.D{{"$set", schema}},
		opts,
	)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	return repo.GetById(ctx, id)
}

func (repo *Repository[DType, SType]) UpdateField(ctx context.Context, id primitive.ObjectID, field string, value any) (bool, *gerror.Error) {
	//obId, _ := primitive.ObjectIDFromHex(id)
	_, err := repo.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set", bson.D{{field, value}}}},
	)
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}
func (repo *Repository[DType, SType]) UpdateFields(ctx context.Context, id primitive.ObjectID, fields map[string]any) (bool, *gerror.Error) {
	//obId, _ := primitive.ObjectIDFromHex(id)

	set := bson.D{}
	for key, value := range fields {
		set = append(set, bson.E{Key: key, Value: value})
	}
	_, err := repo.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set", set}},
	)
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}

func (repo *Repository[DType, SType]) Delete(ctx context.Context, id primitive.ObjectID) (bool, *gerror.Error) {
	//obId, _ := primitive.ObjectIDFromHex(id)
	_, err := repo.Collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}

func (repo *Repository[DType, SType]) Count(ctx context.Context, query bson.D) (int64, *gerror.Error) {
	opts := options.Count().SetHint("_id_")
	count, err := repo.Collection.CountDocuments(ctx, query, opts)
	if err != nil {
		return 0, gerror.GetError(err)
	}
	return count, nil
}
