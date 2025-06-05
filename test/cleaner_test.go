package test

import (
	"bewell_test/models"
	"bewell_test/services"
	"reflect"
	"testing"
)

func TestCleaner_Case1(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         50,
			TotalPrice:        100,
		},
	}

	got := services.CleanOrders(input)

	if len(got) != 3 {
		t.Errorf("expected 3 items, got %d", len(got))
	}
	if got[0].ProductId != "FG0A-CLEAR-IPHONE16PROMAX" {
		t.Errorf("unexpected main productId: %v", got[0].ProductId)
	}
	if got[1].ProductId != "WIPING-CLOTH" {
		t.Errorf("expected WIPING-CLOTH")
	}
	if got[2].ProductId != "CLEAR-CLEANNER" {
		t.Errorf("expected CLEAR-CLEANNER")
	}
}

func TestCleaner_Case2(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         50,
			TotalPrice:        100,
		},
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Qty:        2,
			UnitPrice:  50,
			TotalPrice: 100,
		},
		{
			No:         2,
			ProductId:  "WIPING-CLOTH",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         3,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	result := services.CleanOrders(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case2: Expected %+v, but got %+v", expected, result)
	}
}

func TestCleaner_Case3(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
			Qty:               1,
			UnitPrice:         90, // ไม่ใช้ตรงนี้แล้ว
			TotalPrice:        90,
		},
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
			MaterialId: "FG0A-MATTE",
			ModelId:    "IPHONE16PROMAX",
			Qty:        3,
			UnitPrice:  30.00,
			TotalPrice: 90.00,
		},
		{
			No:         2,
			ProductId:  "WIPING-CLOTH",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         3,
			ProductId:  "MATTE-CLEANNER",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	result := services.CleanOrders(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Case3: Expected %+v, but got %+v", expected, result)
	}
}

func TestCleanOrders_Case4(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
			Qty:               1,
			UnitPrice:         80,
			TotalPrice:        80,
		},
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         2,
			ProductId:  "FG0A-CLEAR-OPPOA3-B",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3-B",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         3,
			ProductId:  "WIPING-CLOTH",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         4,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	got := services.CleanOrders(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Case4: expected %+v, got %+v", expected, got)
	}
}
func TestCleanOrders_Case5(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
			Qty:               1,
			UnitPrice:         120,
			TotalPrice:        120,
		},
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         2,
			ProductId:  "FG0A-CLEAR-OPPOA3-B",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3-B",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         3,
			ProductId:  "FG0A-MATTE-OPPOA3",
			MaterialId: "FG0A-MATTE",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         4,
			ProductId:  "WIPING-CLOTH",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         5,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         6,
			ProductId:  "MATTE-CLEANNER",
			Qty:        1,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	got := services.CleanOrders(input)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Case5: expected %+v, but got %+v", expected, got)
	}
}

func TestCleanOrders_Case6(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
			Qty:               1,
			UnitPrice:         120,
			TotalPrice:        120,
		},
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3",
			Qty:        2,
			UnitPrice:  40,
			TotalPrice: 80,
		},
		{
			No:         2,
			ProductId:  "FG0A-MATTE-OPPOA3",
			MaterialId: "FG0A-MATTE",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         3,
			ProductId:  "WIPING-CLOTH",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         4,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         5,
			ProductId:  "MATTE-CLEANNER",
			Qty:        1,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	got := services.CleanOrders(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Case6: expected %+v, got %+v", expected, got)
	}
}

func TestCleanOrders_Case7(t *testing.T) {
	input := []models.InputOrder{
		{
			No:                1,
			PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
			Qty:               1,
			UnitPrice:         160,
			TotalPrice:        160,
		},
		{
			No:                2,
			PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
			Qty:               1,
			UnitPrice:         50,
			TotalPrice:        50,
		},
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3",
			Qty:        2,
			UnitPrice:  40,
			TotalPrice: 80,
		},
		{
			No:         2,
			ProductId:  "FG0A-MATTE-OPPOA3",
			MaterialId: "FG0A-MATTE",
			ModelId:    "OPPOA3",
			Qty:        2,
			UnitPrice:  40,
			TotalPrice: 80,
		},
		{
			No:         3,
			ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
			MaterialId: "FG0A-PRIVACY",
			ModelId:    "IPHONE16PROMAX",
			Qty:        1,
			UnitPrice:  50,
			TotalPrice: 50,
		},
		{
			No:         4,
			ProductId:  "WIPING-CLOTH",
			Qty:        5,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         5,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         6,
			ProductId:  "MATTE-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         7,
			ProductId:  "PRIVACY-CLEANNER",
			Qty:        1,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	got := services.CleanOrders(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Case7: expected %+v, got %+v", expected, got)
	}
}
