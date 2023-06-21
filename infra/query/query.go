package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PBKKE08/FP-BE/api/query"
	"github.com/PBKKE08/FP-BE/core/model/partner"
	"github.com/PBKKE08/FP-BE/core/model/pengguna"
	"github.com/jmoiron/sqlx"
)

type Query struct {
	db *sqlx.DB
}

func NewQuery(db *sqlx.DB) *Query {
	return &Query{db: db}
}

func (q *Query) By(ctx context.Context, daerah string, jenisKelamin string, kebutuhan string) []query.CariPasangan {
	var results []query.CariPasangan

	switch {
	case daerah != "" && jenisKelamin != "" && kebutuhan != "":
		qr := `
		SELECT
		  p.id,
		  p.name,
		  p.price,
		  p.gender,
		  c.name as c_name,
		  ca.name as cat_name,
		  COALESCE(average_rating, 0) as rating
		FROM
		  partners p
		  JOIN cities c ON p.city_id = c.id
		  JOIN categories ca on p.category_id = ca.id
		  LEFT JOIN (
			SELECT
			  partner_id,
			  AVG(rating) AS average_rating
			FROM
			  reviews
			GROUP BY
			  partner_id
		  ) r ON p.id = r.partner_id
		WHERE c.name = ? AND p.gender = ? AND ca.name = ?;
		`

		rows, _ := q.db.QueryxContext(ctx, qr, daerah, jenisKelamin, kebutuhan)
		for rows.Next() {
			var qq query.CariPasangan
			rows.StructScan(&qq)

			results = append(results, qq)
		}

		return results

	case daerah != "":
		qr := `
		SELECT
		  p.id,
		  p.name,
		  p.price,
		  p.gender,
		  c.name as c_name,
		  ca.name as cat_name,
		  COALESCE(average_rating, 0) as rating
		FROM
		  partners p
		  JOIN cities c ON p.city_id = c.id
		  JOIN categories ca on p.category_id = ca.id
		  LEFT JOIN (
			SELECT
			  partner_id,
			  AVG(rating) AS average_rating
			FROM
			  reviews
			GROUP BY
			  partner_id
		  ) r ON p.id = r.partner_id
		WHERE c.name = ?;
		`

		rows, _ := q.db.QueryxContext(ctx, qr, daerah)
		for rows.Next() {
			var qq query.CariPasangan
			rows.StructScan(&qq)

			results = append(results, qq)
		}

		return results

	case jenisKelamin != "":
		qr := `
		SELECT
		  p.id,
		  p.name,
		  p.price,
		  p.gender,
		  c.name as c_name,
		  ca.name as cat_name,
		  COALESCE(average_rating, 0) as rating
		FROM
		  partners p
		  JOIN cities c ON p.city_id = c.id
		  JOIN categories ca on p.category_id = ca.id
		  LEFT JOIN (
			SELECT
			  partner_id,
			  AVG(rating) AS average_rating
			FROM
			  reviews
			GROUP BY
			  partner_id
		  ) r ON p.id = r.partner_id
		WHERE p.gender = ?;
		`

		rows, _ := q.db.QueryxContext(ctx, qr, jenisKelamin)
		for rows.Next() {
			var qq query.CariPasangan
			rows.StructScan(&qq)

			results = append(results, qq)
		}

		return results

	case kebutuhan != "":
		qr := `
		SELECT
		  p.id,
		  p.name,
		  p.price,
		  p.gender,
		  c.name as c_name,
		  ca.name as cat_name,
		  COALESCE(average_rating, 0) as rating
		FROM
		  partners p
		  JOIN cities c ON p.city_id = c.id
		  JOIN categories ca on p.category_id = ca.id
		  LEFT JOIN (
			SELECT
			  partner_id,
			  AVG(rating) AS average_rating
			FROM
			  reviews
			GROUP BY
			  partner_id
		  ) r ON p.id = r.partner_id
		WHERE ca.name = ?;`

		rows, _ := q.db.QueryxContext(ctx, qr, kebutuhan)
		for rows.Next() {
			var qq query.CariPasangan
			rows.StructScan(&qq)

			results = append(results, qq)
		}

		return results

	default:
		qr := `
		SELECT
		  p.id,
		  p.name,
		  p.price,
		  p.gender,
		  c.name as c_name,
		  ca.name as cat_name,
		  COALESCE(average_rating, 0) as rating
		FROM
		  partners p
		  JOIN categories ca on p.category_id = ca.id
		  JOIN cities c ON p.city_id = c.id
		  LEFT JOIN (
			SELECT
			  partner_id,
			  AVG(rating) AS average_rating
			FROM
			  reviews
			GROUP BY
			  partner_id
		  ) r ON p.id = r.partner_id;
		`

		rows, err := q.db.QueryxContext(ctx, qr)
		if err != nil {
			println(err.Error())
		}

		for rows.Next() {
			var qq query.CariPasangan
			errs := rows.StructScan(&qq)
			if errs != nil {
				println(errs.Error())
			}

			results = append(results, qq)
		}

		return results
	}
}

func (q *Query) ByUserEmail(ctx context.Context, email string) query.CariUserByEmail {
	var results query.CariUserByEmail

	qr := `SELECT id, name, email, telephone, gender FROM users WHERE email = ?`

	q.db.GetContext(ctx, &results, qr, email)

	return results
}

func (q *Query) GetAllCityAndCategory(ctx context.Context) query.AllCitiesAndCategories {
	qr1 := `SELECT * FROM categories`
	qr2 := `SELECT * FROM cities`

	var categories []query.Category
	var cities []query.City

	rows, _ := q.db.QueryxContext(ctx, qr1)
	defer rows.Close()

	for rows.Next() {
		var c query.Category

		rows.StructScan(&c)

		categories = append(categories, c)
	}

	rowsCity, _ := q.db.QueryxContext(ctx, qr2)
	defer rowsCity.Close()

	for rowsCity.Next() {
		var c query.City

		rowsCity.StructScan(&c)

		cities = append(cities, c)
	}

	var all query.AllCitiesAndCategories

	all.Categories = categories
	all.Cities = cities

	return all
}

func (q *Query) LihatPartnerDetail(ctx context.Context, id partner.ID) (query.DetailPartner, error) {
	qr := `
		SELECT
		  p.id as id,
		  p.email as email,
		  p.name as name,
		  p.price as price,
		  p.gender as gender,
		  p.description as description,
		  c.name as city_name,
		  ca.name as cat_name,
		  COALESCE(average_rating, 0) as rating
		FROM
		  partners p
		  JOIN categories ca on p.category_id = ca.id
		  JOIN cities c ON p.city_id = c.id
		  LEFT JOIN (
			SELECT
			  partner_id,
			  AVG(rating) AS average_rating
			FROM
			  reviews
			GROUP BY
			  partner_id
		  ) r ON p.id = r.partner_id
		WHERE p.id = ?;
		`

	var detailPartner query.DetailPartner

	err := q.db.GetContext(ctx, &detailPartner, qr, id.String())

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return detailPartner, fmt.Errorf("error getting partners: %w", err)
		}

		return detailPartner, fmt.Errorf("partner not found")
	}

	return detailPartner, nil
}

func (q *Query) LihatTransaksi(ctx context.Context, id pengguna.ID) ([]query.SeluruhTransaksi, error) {
	var results []query.SeluruhTransaksi

	qr := `
		SELECT 
		p.name, 
		c.name as cat_name, 
		o.booking_day as booking_date, 
		o.time_start as start_time, 
		o.time_end as end_time,
		CASE 
			WHEN o.booking_day < curdate() THEN 1
			WHEN o.booking_day = curdate() AND o.time_end > curtime() THEN 1
			WHEN o.booking_day = curdate() AND o.time_start > curtime() AND o.time_end < curtime() THEN 2
			WHEN o.booking_day = curdate() AND o.time_start < curtime() THEN 3
			ELSE 3
		END AS order_status
		
		FROM partners p
		
		JOIN categories c
		ON p.category_id = c.id
		
		JOIN orders o
		ON p.id = o.partners_id
		
		WHERE o.user_id = ?;
		`

	rows, err := q.db.QueryxContext(ctx, qr, id.String())
	if err != nil {
		return results, err
	}

	defer rows.Close()

	for rows.Next() {
		var qq query.SeluruhTransaksi

		if err := rows.StructScan(&qq); err != nil {
			return results, err
		}

		results = append(results, qq)
	}

	return results, err
}
