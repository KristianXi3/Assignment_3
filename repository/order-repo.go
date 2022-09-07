package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/KristianXi3/Assignment_3/entity"
)

type OrderRepoIface interface {
	GetOrders(ctx context.Context) ([]*entity.Order, error)
	GetOrderById(ctx context.Context, id int) (*entity.Order, error)
	CreateOrder(ctx context.Context, order entity.Order) (string, error)
	UpdateOrder(ctx context.Context, id int, order entity.Order) (string, error)
	DeleteOrder(ctx context.Context, id int) (string, error)
}

type OrderRepo struct {
	sql *sql.DB
}

func NewOrderRepo(context *sql.DB) OrderRepoIface {
	return &OrderRepo{sql: context}
}

func (repo *OrderRepo) GetOrders(ctx context.Context) ([]*entity.Order, error) {
	orders := []*entity.Order{}

	data, err := repo.sql.QueryContext(ctx, "SELECT order_id, customer_name, ordered_at FROM Orders")
	defer data.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for data.Next() {
		var order entity.Order
		err := data.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	var wg sync.WaitGroup
	wg.Add(len(orders))
	for i := 0; i < len(orders); i++ {
		go func(x int) {
			defer wg.Done()
			data, err = repo.sql.QueryContext(ctx, "SELECT item_id, item_code, description, quantity FROM Items WHERE order_id = @Order_Id ",
				sql.Named("Order_Id", orders[x].OrderId))
			defer data.Close()
			if err != nil {
				log.Fatal(err)
			}
			items := []entity.Item{}
			for data.Next() {
				var item entity.Item
				err := data.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity)
				if err != nil {
					log.Fatal(err)
				}
				items = append(items, item)
			}
			orders[x].Items = items
			//items = []entity.Item{}
		}(i)
	}
	wg.Wait()
	return orders, nil
}

func (repo *OrderRepo) GetOrderById(ctx context.Context, id int) (*entity.Order, error) {
	var order entity.Order

	data, err := repo.sql.QueryContext(ctx, "SELECT order_id, customer_name, ordered_at FROM Orders WHERE order_id = @Id",
		sql.Named("Id", id))
	defer data.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for data.Next() {
		err := data.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	items := []entity.Item{}
	data, err = repo.sql.QueryContext(ctx, "SELECT item_Id,  item_code, description, quantity FROM Items WHERE order_id = @Id",
		sql.Named("Id", id))
	defer data.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for data.Next() {
		var item entity.Item
		err := data.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

func (repo *OrderRepo) CreateOrder(ctx context.Context, order entity.Order) (string, error) {
	var result string

	data, err := repo.sql.QueryContext(ctx, "INSERT into ORDERS (customer_name, ordered_at) values ( @customer_name, @ordered_at); select order_id = convert(bigint, SCOPE_IDENTITY())",
		sql.Named("customer_name", order.CustomerName),
		sql.Named("ordered_at", order.OrderedAt))
	if err != nil {
		log.Fatal(err)
		return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
	}
	defer data.Close()

	var lastOrderId int
	for data.Next() {
		err := data.Scan(&lastOrderId)
		if err != nil {
			log.Fatal(err)
			return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
		}
	}

	for i := 0; i < len(order.Items); i++ {
		_, err = repo.sql.ExecContext(ctx, "INSERT into ITEMS (item_code, description,quantity, order_id) values (@code, @description, @quantity, @order_id)",
			sql.Named("code", order.Items[i].ItemCode),
			sql.Named("description", order.Items[i].Description),
			sql.Named("quantity", order.Items[i].Quantity),
			sql.Named("order_id", lastOrderId))

		if err != nil {
			log.Fatal(err)
			return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
		}
	}

	result = "Order created successfully"
	return result, nil
}

func (repo *OrderRepo) UpdateOrder(ctx context.Context, id int, order entity.Order) (string, error) {
	result := ""

	_, err := repo.sql.ExecContext(ctx, "UPDATE ORDERS set customer_name = @customer_name, ordered_at = @ordered_at where order_id = @id",
		sql.Named("customer_name", order.CustomerName),
		sql.Named("ordered_at", order.OrderedAt),
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
		return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
	}

	for i := 0; i < len(order.Items); i++ {
		var item entity.Item
		exist := false
		data, err := repo.sql.QueryContext(ctx, "SELECT item_id,  item_code, description, quantity FROM Items WHERE order_id = @id AND item_id = @itemId",
			sql.Named("id", id),
			sql.Named("itemId", order.Items[i].ItemId))

		if err != nil {
			log.Fatal(err)
			return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
		}
		defer data.Close()

		for data.Next() {
			err := data.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity)
			if err != nil {
				log.Fatal(err)
				return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
			}
		}

		if (item != entity.Item{}) {
			_, err := repo.sql.ExecContext(ctx, "UPDATE ITEMS set item_code = @item_code, description = @description, quantity = @quantity WHERE item_id = @item_Id",
				sql.Named("item_code", order.Items[i].ItemCode),
				sql.Named("description", order.Items[i].Description),
				sql.Named("quantity", order.Items[i].Quantity),
				sql.Named("item_Id", item.ItemId))

			if err != nil {
				log.Fatal(err)
				return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
			}
			exist = true
		}

		if !exist {
			result += fmt.Sprintf("Item id %d not found \n", order.Items[i].ItemId)
		}
	}
	result += "Updated successfully"
	return result, nil
}

func (repo *OrderRepo) DeleteOrder(ctx context.Context, id int) (string, error) {
	var result string

	_, err := repo.sql.ExecContext(ctx, "DELETE from items where order_id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
		return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
	}

	_, err = repo.sql.ExecContext(ctx, "DELETE from ORDERS where order_id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
		return fmt.Sprintf("Internal Server Error: %s", err.Error()), err
	}

	result = fmt.Sprintf("Order %d deleted", id)
	return result, nil
}
