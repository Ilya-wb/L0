package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Delivery struct {
	Delivery_name string
	Phone         string
	Zip           string
	City          string
	Address       string
	Region        string
	Email         string
}

type Payment struct {
	Transaction  string
	RequestId    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFree   int
}

type Item struct {
	ChrtId     int
	TrackItem  string
	Price      int
	Rid        string
	Item_name  string
	Sale       int
	Item_size  string
	TotalPrice int
	NmId       int
	Brand      string
	Status     int
}

type Order struct {
	OrderUid          string
	TrackNumber       string
	Entry             string
	Delivery          Delivery
	Payment           Payment
	Items             [1]Item
	Locale            string
	InternalSignature string
	CustomerId        string
	DeliveryService   string
	Shardkey          string
	SmId              int
	DateCreated       string
	OofShard          string
}

var database *sql.DB
var chachOrders []Order

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		form := r.FormValue("id")

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		fmt.Println("Form before for 1", form)
		if form != "" {
			for _, val := range chachOrders {
				if val.OrderUid == form {
					//fmt.Println("contact", form)
					Json_data := val
					st, _ := json.Marshal(Json_data)
					_, err := fmt.Fprintln(w, string(st))
					if err != nil {
						fmt.Println(err)
					}
					break
				}
			}
		}
	})

	connStr := "user=postgres password=raze1998 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Println(err)
	}

	database = db

	defer db.Close()

	//tmpl := template.Must(template.ParseFiles("index.html"))

	rows, err := database.Query("select * from allOrders")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard, &order.Delivery.Delivery_name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFree, &order.Items[0].ChrtId, &order.Items[0].Price, &order.Items[0].Rid, &order.Items[0].Item_name, &order.Items[0].Sale, &order.Items[0].Item_size, &order.Items[0].TotalPrice, &order.Items[0].NmId, &order.Items[0].Brand, &order.Items[0].Status, &order.Items[0].TrackItem)
		if err != nil {
			fmt.Println(err)
		} else {
			chachOrders = append(chachOrders, order)
		}
	}

	fmt.Print("//////////\nКэш, восстановленый из БД\n//////////\n")
	for i, _ := range chachOrders {
		fmt.Println("Элемент ", i)
		fmt.Println(chachOrders[i])
	}

	type Val struct {
		id    string
		value string
	}

	var ch Order
	sc, _ := stan.Connect("test-cluster", "myId")
	sc.Subscribe("postgres", func(m *stan.Msg) {
		mess := strings.ToLower(string(m.Data))
		str := strings.Split(mess, ",")
		var arr []string
		for _, val := range str {
			arr = strings.Split(val, "-")
			switch arr[0] {
			case "orderuid":
				ch.OrderUid = arr[1]
			case "tracknumber":
				ch.TrackNumber = arr[1]
			case "entry":
				ch.Entry = arr[1]
			case "locale":
				ch.Locale = arr[1]
			case "internalsignature":
				ch.InternalSignature = arr[1]
			case "customerId":
				ch.CustomerId = arr[1]
			case "deliveryservice":
				ch.DeliveryService = arr[1]
			case "shardkey":
				ch.Shardkey = arr[1]
			case "smid":
				ch.SmId, err = strconv.Atoi(arr[1])
			case "datecreated":
				ch.DateCreated = arr[1]
			case "status":
				ch.OofShard = arr[1]
			case "delivery_name":
				//
				ch.Delivery.Delivery_name = arr[1]
			case "phone":
				ch.Delivery.Phone = arr[1]
			case "zip":
				ch.Delivery.Zip = arr[1]
			case "city":
				ch.Delivery.City = arr[1]
			case "address":
				ch.Delivery.Address = arr[1]
			case "region":
				ch.Delivery.Region = arr[1]
			case "email":
				ch.Delivery.Email = arr[1]
			case "transaction":
				//
				ch.Payment.Transaction = arr[1]
			case "requestid":
				ch.Payment.RequestId = arr[1]
			case "currency":
				ch.Payment.Currency = arr[1]
			case "provider":
				ch.Payment.Provider = arr[1]
			case "amount":
				ch.Payment.Amount, err = strconv.Atoi(arr[1])
			case "paymentDt":
				ch.Payment.PaymentDt, err = strconv.Atoi(arr[1])
			case "bank":
				ch.Payment.Bank = arr[1]
			case "deliverycost":
				ch.Payment.DeliveryCost, err = strconv.Atoi(arr[1])
			case "goodstotal":
				ch.Payment.GoodsTotal, err = strconv.Atoi(arr[1])
			case "customfree":
				//
				ch.Payment.CustomFree, err = strconv.Atoi(arr[1])
			case "chrtid":
				ch.Items[0].ChrtId, err = strconv.Atoi(arr[1])
			case "trackitem":
				ch.Items[0].TrackItem = arr[1]
			case "price":
				ch.Items[0].Price, err = strconv.Atoi(arr[1])
			case "rid":
				ch.Items[0].Rid = arr[1]
			case "item_name":
				ch.Items[0].Item_name = arr[1]
			case "sale":
				ch.Items[0].Sale, err = strconv.Atoi(arr[1])
			case "item_size":
				ch.Items[0].Item_size = arr[1]
			case "totalprice":
				ch.Items[0].TotalPrice, err = strconv.Atoi(arr[1])
			case "nmid":
				ch.Items[0].NmId, err = strconv.Atoi(arr[1])
			case "brand":
				ch.Items[0].Brand = arr[1]
			case "oofshard":
				ch.Items[0].Status, err = strconv.Atoi(arr[1])
				rows, err = database.Query("insert into allOrders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_name, phone, zip, city, address, region, email, transaction_, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, chrt_id, price, rid, item_name, sale, item_size, total_price, nm_id, brand, status, track_item)"+
					" values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39);", ch.OrderUid, ch.TrackNumber, ch.Entry, ch.Locale, ch.InternalSignature, ch.CustomerId, ch.DeliveryService, ch.Shardkey, ch.SmId, ch.DateCreated, ch.OofShard, ch.Delivery.Delivery_name, ch.Delivery.Phone, ch.Delivery.Zip, ch.Delivery.City, ch.Delivery.Address, ch.Delivery.Region, ch.Delivery.Email, ch.Payment.Transaction, ch.Payment.RequestId, ch.Payment.Currency, ch.Payment.Provider, ch.Payment.Amount, ch.Payment.PaymentDt, ch.Payment.Bank, ch.Payment.DeliveryCost, ch.Payment.GoodsTotal, ch.Payment.CustomFree, ch.Items[0].ChrtId, ch.Items[0].Price, ch.Items[0].Rid, ch.Items[0].Item_name, ch.Items[0].Sale, ch.Items[0].Item_size, ch.Items[0].TotalPrice, ch.Items[0].NmId, ch.Items[0].Brand, ch.Items[0].Status, ch.Items[0].TrackItem)
				fmt.Println("добавленый элемент")
				fmt.Println(ch)
			}
		}
		chachOrders = append(chachOrders, ch)
		//if err != nil {
		//fmt.Println(err)
		//fmt.Println(rows)
		//}
	}, stan.DurableName("my-durable"))

	fmt.Println("Server is listening...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
