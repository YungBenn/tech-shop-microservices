package repository

import (
	"context"
	"time"

	"github.com/YungBenn/tech-shop-microservices/internal/mongodb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	InsertProduct(ctx context.Context, product entity.Product) (*entity.Product, error)
	FindAllProducts(ctx context.Context, limit int, page int) ([]entity.Product, *mongodb.MongoPaginate, error)
	FindOneProduct(ctx context.Context, id string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id string, product entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type ProductRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
	log  *logrus.Logger
}

func NewProductRepository(db *mongo.Database, log *logrus.Logger) ProductRepository {
	return &ProductRepositoryImpl{
		db:   db,
		coll: db.Collection("products"),
		log:  log,
	}
}

func (p *ProductRepositoryImpl) DeleteProduct(ctx context.Context, id string) error {
	coll := p.coll

	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		p.log.Error("Error deleting product: ", err)
		return err
	}

	return nil
}

func (p *ProductRepositoryImpl) FindAllProducts(ctx context.Context, limit int, page int) ([]entity.Product, *mongodb.MongoPaginate, error) {
	coll := p.coll

	pagination := mongodb.NewMongoPaginate(limit, page)

	totalRows, err := coll.CountDocuments(ctx, bson.M{})
    if err != nil {
        p.log.Error("Error counting products: ", err)
        return nil, nil, err
    }

	totalPages := pagination.CalculateTotalPages(totalRows)

	paginationData := mongodb.MongoPaginate{
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	res, err := coll.Find(ctx, bson.M{}, pagination.GetPaginatedOpts())
	if err != nil {
		p.log.Error("Error finding products: ", err)
		return nil, nil, err
	}

	products := make([]entity.Product, 0)
	if err = res.All(ctx, &products); err != nil {
		p.log.Error("Error decoding products: ", err)
		return nil, nil, err
	}

	return products, &paginationData, nil
}

func (p *ProductRepositoryImpl) FindOneProduct(ctx context.Context, id string) (*entity.Product, error) {
	coll := p.coll

	objectId, _ := primitive.ObjectIDFromHex(id)

	var product entity.Product
	err := coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&product)
	if err != nil {
		p.log.Error("Error finding product: ", err)
		return nil, err
	}

	return &product, nil
}

func (p *ProductRepositoryImpl) InsertProduct(ctx context.Context, product entity.Product) (*entity.Product, error) {
	coll := p.coll

	res, err := coll.InsertOne(ctx, product)
	if err != nil {
		p.log.Error("Error inserting product: ", err)
		return nil, err
	}

	newID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		p.log.Error("Error asserting InsertedID to primitive.ObjectID")
		return nil, err
	}

	newProduct := &entity.Product{
		ID:          newID,
		Title:       product.Title,
		Price:       product.Price,
		Tag:         product.Tag,
		Discount:    product.Discount,
		Stock:       product.Stock,
		Image:       product.Image,
		Description: product.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return newProduct, nil
}

func (p *ProductRepositoryImpl) UpdateProduct(ctx context.Context, id string, product entity.Product) (*entity.Product, error) {
	coll := p.coll

	objectId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": product}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		p.log.Error("Error updating product: ", err)
		return nil, err
	}

	var updatedProduct entity.Product
	err = coll.FindOne(ctx, filter).Decode(&updatedProduct)
	if err != nil {
		p.log.Error("Error finding updated product: ", err)
		return nil, err
	}

	return &updatedProduct, nil
}
