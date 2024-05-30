package main

import "fmt"

type Product struct {
	Id          int
	Name        string
	Description string
	Icon        string
}

type Store struct {
	products map[int]Product
}

func (s *Store) AddProduct(name string, description string, iconPath string, currentId int) Product {
	currentId += 1
	product := Product{currentId, name, description, iconPath}
	s.products[currentId] = product
	return product
}

func (s *Store) idChecker(id int) error {
	if id < 1 {
		return fmt.Errorf("incorrect product id: %d", id)
	}

	if id != s.products[id].Id {
		return fmt.Errorf("no such product. Given id: %d", id)
	}

	return nil
}

func (s *Store) GetProduct(id int) (Product, error) {
	var err = s.idChecker(id)
	if err != nil {
		return Product{}, err
	}

	return s.products[id], nil
}

func (s *Store) DeleteProduct(id int) (Product, error) {
	var err = s.idChecker(id)
	if err != nil {
		return Product{}, err
	}

	var product = s.products[id]
	delete(s.products, id)

	// defer delete(s.products, id) maybe like this?

	return product, nil
}

func (s *Store) UpdateProductData(id int, field string, value string) (Product, error) {
	var err = s.idChecker(id)
	if err != nil {
		return Product{}, err
	}

	var product = s.products[id]

	switch field {
	case "name":
		product.Name = value
	case "description":
		product.Description = value
	case "icon":
		product.Icon = value
	}
	s.products[id] = product

	return product, nil
}

func (s *Store) GetAllProducts() []Product {
	productList := []Product{}

	for _, v := range s.products {
		productList = append(productList, v)
	}

	return productList
}

func (s *Store) AddProductIcon(id int, iconPath string) (Product, error) {
	var err = s.idChecker(id)
	if err != nil {
		return Product{}, err
	}

	var productCopy = s.products[id]
	productCopy.Icon = iconPath
	s.products[id] = productCopy

	return s.products[id], nil
}

func (s *Store) GetProductIcon(id int) (Product, error) {
	var err = s.idChecker(id)
	if err != nil {
		return Product{}, err
	}

	return s.products[id], nil
}
