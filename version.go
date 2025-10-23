package tea

import (
	"github.com/bytedance/sonic"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
)

const Version = "v0.0.0"

func A() {
	mysql.NewConfig()
	_, _ = pgx.ParseConfig("")
	_ = sonic.Unmarshal(nil, nil)
}

// Inserting Records

// INSERT INTO t1 VALUES (1), (2), (3);
// INSERT INTO t2 VALUES (2), (4);

// INSERT INTO student_tests
//  (name, test, score, test_date) VALUES
//  ('Chun', 'SQL', 75, '2012-11-05'),
//  ('Chun', 'Tuning', 73, '2013-06-14'),
//  ('Esben', 'SQL', 43, '2014-02-11'),
//  ('Esben', 'Tuning', 31, '2014-02-09'),
//  ('Kaolin', 'SQL', 56, '2014-01-01'),
//  ('Kaolin', 'Tuning', 88, '2013-12-29'),
//  ('Tatiana', 'SQL', 87, '2012-04-28'),
//  ('Tatiana', 'Tuning', 83, '2013-09-30');

// Querying from two tables on a common value

// SELECT * FROM t1 INNER JOIN t2 ON t1.a = t2.b;

// Ordering Results

// SELECT name, test, score FROM student_tests ORDER BY score DESC;
// +---------+--------+-------+
// | name    | test   | score |
// +---------+--------+-------+
// | Kaolin  | Tuning |    88 |
// | Tatiana | SQL    |    87 |
// | Tatiana | Tuning |    83 |
// | Chun    | SQL    |    75 |
// | Chun    | Tuning |    73 |
// | Kaolin  | SQL    |    56 |
// | Esben   | SQL    |    43 |
// | Esben   | Tuning |    31 |
// +---------+--------+-------+
