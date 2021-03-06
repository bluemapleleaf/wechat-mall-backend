package dbops

import (
	"strconv"
	"time"
	"wechat-mall-backend/model"
)

const skuColumnList = `
id, title, price, code, stock, goods_id, online, picture, specs, is_del, create_time, update_time
`

func GetSKUList(title string, goodsId, online, page, size int) (*[]model.WechatMallSkuDO, error) {
	sql := "SELECT " + skuColumnList + " FROM wechat_mall_sku WHERE is_del = 0"
	if goodsId != 0 {
		sql += " AND goods_id = " + strconv.Itoa(goodsId)
	}
	if title != "" {
		sql += " AND title like '%" + title + "%'"
	}
	if online == 0 || online == 1 {
		sql += " AND online = " + strconv.Itoa(online)
	}
	if page > 0 && size > 0 {
		sql += " LIMIT " + strconv.Itoa((page-1)*size) + ", " + strconv.Itoa(size)
	}
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	var skuList []model.WechatMallSkuDO
	for rows.Next() {
		sku := model.WechatMallSkuDO{}
		err := rows.Scan(&sku.Id, &sku.Title, &sku.Price, &sku.Code, &sku.Stock, &sku.GoodsId, &sku.Online, &sku.Picture,
			&sku.Specs, &sku.Del, &sku.CreateTime, &sku.UpdateTime)
		if err != nil {
			return nil, err
		}
		skuList = append(skuList, sku)
	}
	return &skuList, nil
}

func CountSKU(title string, goodsId, online int) (int, error) {
	sql := "SELECT COUNT(*) FROM wechat_mall_sku WHERE is_del = 0"
	if goodsId != 0 {
		sql += " AND goods_id = " + strconv.Itoa(goodsId)
	}
	if online == 0 || online == 1 {
		sql += " AND online = " + strconv.Itoa(online)
	}
	if title != "" {
		sql += " AND title like '%" + title + "%'"
	}
	rows, err := dbConn.Query(sql)
	if err != nil {
		return 0, err
	}
	total := 0
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}
	return total, nil
}

func AddSKU(sku *model.WechatMallSkuDO) (int64, error) {
	sql := "INSERT INTO wechat_mall_sku( " + skuColumnList[4:] + " ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(sku.Title, sku.Price, sku.Code, sku.Stock, sku.GoodsId, sku.Online, sku.Picture, sku.Specs, 0, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetSKUById(id int) (*model.WechatMallSkuDO, error) {
	sql := "SELECT " + skuColumnList + " FROM wechat_mall_sku WHERE id = " + strconv.Itoa(id)
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	sku := model.WechatMallSkuDO{}
	for rows.Next() {
		err := rows.Scan(&sku.Id, &sku.Title, &sku.Price, &sku.Code, &sku.Stock, &sku.GoodsId, &sku.Online, &sku.Picture,
			&sku.Specs, &sku.Del, &sku.CreateTime, &sku.UpdateTime)
		if err != nil {
			return nil, err
		}
	}
	return &sku, nil
}

func GetSKUByCode(code string) (*model.WechatMallSkuDO, error) {
	sql := "SELECT " + skuColumnList + " FROM wechat_mall_sku WHERE is_del = 0 AND code = '" + code + "'"
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	sku := model.WechatMallSkuDO{}
	for rows.Next() {
		err := rows.Scan(&sku.Id, &sku.Title, &sku.Price, &sku.Code, &sku.Stock, &sku.GoodsId, &sku.Online, &sku.Picture,
			&sku.Specs, &sku.Del, &sku.CreateTime, &sku.UpdateTime)
		if err != nil {
			return nil, err
		}
	}
	return &sku, nil
}

func UpdateSKUById(sku *model.WechatMallSkuDO) error {
	sql := `
UPDATE wechat_mall_sku
SET title = ?, price = ?, code = ?, stock = ?, goods_id = ?, online = ?, picture = ?, specs = ?, is_del = ?, update_time = ?
WHERE id = ?
`
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sku.Title, sku.Price, sku.Code, sku.Stock, sku.GoodsId, sku.Online, sku.Picture, sku.Specs, sku.Del, time.Now(), sku.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSkuStockById(id, num int) error {
	sql := "UPDATE wechat_mall_sku SET update_time = now(), stock = stock - " + strconv.Itoa(num) + " WHERE id = " + strconv.Itoa(id)
	sql += " AND stock >= " + strconv.Itoa(num)
	_, err := dbConn.Exec(sql)
	return err
}

func QuerySellOutSKUList(page, size int) (*[]model.WechatMallSkuDO, error) {
	sql := "SELECT " + skuColumnList + " FROM wechat_mall_sku WHERE is_del = 0 AND stock = 0"
	if page > 0 && size > 0 {
		sql += " LIMIT " + strconv.Itoa(page) + ", " + strconv.Itoa(size)
	}
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	skuList := []model.WechatMallSkuDO{}
	for rows.Next() {
		sku := model.WechatMallSkuDO{}
		err := rows.Scan(&sku.Id, &sku.Title, &sku.Price, &sku.Code, &sku.Stock, &sku.GoodsId, &sku.Online, &sku.Picture,
			&sku.Specs, &sku.Del, &sku.CreateTime, &sku.UpdateTime)
		if err != nil {
			return nil, err
		}
		skuList = append(skuList, sku)
	}
	return &skuList, nil
}

func CountSellOutSKUList() (int, error) {
	sql := "SELECT COUNT(*) FROM wechat_mall_sku WHERE is_del = 0 AND stock = 0"
	rows, err := dbConn.Query(sql)
	if err != nil {
		return 0, err
	}
	total := 0
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}
	return total, nil
}
