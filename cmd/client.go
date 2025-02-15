package main

//import (
//	"bytes"
//	"fmt"
//	"io"
//	"net/http"
//)

//func main() {
//text := "{\n  \"token\": \"testtoken\",\n  \"sku\": 2956315\n}"
//url := "http://route256.pavl.uk:8080/get_product"
//
//jsonBody := []byte(text)
//req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
//if err != nil {
//	fmt.Println(err)
//}
//
//req.Header.Set("Content-Type", "application/json")
//
//client := &http.Client{}
//
//resp, err := client.Do(req)
//if err != nil {
//	fmt.Println(err)
//}
//defer resp.Body.Close()
//
//bufer, err := io.ReadAll(resp.Body)
//if err != nil {
//	panic(err)
//}
//fmt.Println(string(delivery.GetProductHandler(1625903)))
//}
