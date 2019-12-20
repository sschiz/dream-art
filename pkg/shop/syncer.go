/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package shop

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sschiz/dream-art/pkg/category"
	"github.com/sschiz/dream-art/pkg/product"
)

type Syncer struct {
	ConnectionString string
	called           bool
}

// Sync synchronizes the shop with the database
func (s *Syncer) Sync(shop *Shop) error {
	db, err := sql.Open("postgres", s.ConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	if !s.called {
		// Get products and categories from DB
		rows, err := db.Query("SELECT * FROM shop.shop.categories")
		if err != nil {
			return err
		}

		for rows.Next() {
			cat := new(category.Category)

			err := rows.Scan(&cat.Name)
			if err != nil {
				return err
			}

			shop.categories = append(shop.categories, cat)
		}

		stmt, err := db.Prepare("SELECT id, name, description, price, photo FROM shop.shop.products WHERE category_name=$1")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, c := range shop.categories {
			rows, err := stmt.Query(c.Name)
			if err != nil {
				return err
			}

			for rows.Next() {
				pr := new(product.Product)
				err := rows.Scan(&pr.ID, &pr.Name, &pr.Description, &pr.Price, &pr.Photo)
				if err != nil {
					return err
				}

				err = c.AppendProduct(pr.Name, pr.Description, pr.Photo, pr.Price)
				if err != nil {
					return err
				}
			}
		}

		// Get admins from DB
		rows, err = db.Query("SELECT * FROM shop.shop.admins")
		if err != nil {
			return err
		}

		for rows.Next() {
			var nickname string
			var chatID int64

			err := rows.Scan(&nickname, &chatID)
			if err != nil {
				return err
			}
			shop.Admins[nickname] = chatID
		}
		s.called = true
	} else {
		_, err := db.Exec("DELETE FROM shop.shop.products; DELETE FROM shop.shop.categories;")
		if err != nil {
			return err
		}

		newCategoryStmt, err := db.Prepare("INSERT INTO shop.shop.categories VALUES ($1)")
		if err != nil {
			return err
		}
		defer newCategoryStmt.Close()

		newProductStmt, err := db.Prepare("INSERT INTO shop.shop.products (name, description, price, photo, category_name) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			return err
		}
		defer newProductStmt.Close()

		for _, c := range shop.categories {
			_, err := newCategoryStmt.Exec(c.Name)
			if err != nil {
				return err
			}

			for _, p := range c.Products() {
				_, err = newProductStmt.Exec(p.Name, p.Description, p.Price, p.Photo, c.Name)
				if err != nil {
					return err
				}
			}
		}

		_, err = db.Exec("DELETE FROM shop.shop.admins")
		if err != nil {
			return err
		}

		newAdminStmt, err := db.Prepare("INSERT INTO shop.shop.admins VALUES ($1, $2)")
		if err != nil {
			return err
		}
		defer newAdminStmt.Close()

		for nick, chatID := range shop.Admins {
			_, err := newAdminStmt.Exec(nick, chatID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
