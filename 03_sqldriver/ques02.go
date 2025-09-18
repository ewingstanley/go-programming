package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

// transferMoney 执行转账事务
func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    // 开始事务
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("开启事务失败: %v", err)
    }
    // 使用 defer 确保在函数返回时，如果事务未提交则回滚（安全措施）
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p) // 重新抛出 panic
        }
    }()

    // 1. 检查转出账户余额是否足够
    var currentBalance float64
    err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", fromID).Scan(&currentBalance)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("查询账户余额失败: %v", err)
    }
    if currentBalance < amount {
        tx.Rollback()
        return fmt.Errorf("账户 %d 余额不足。当前余额: %.2f, 所需金额: %.2f", fromID, currentBalance, amount)
    }

    // 2. 从转出账户扣除金额
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("扣除转出账户金额失败: %v", err)
    }

    // 3. 向转入账户增加金额
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("增加转入账户金额失败: %v", err)
    }

    // 4. 记录转账交易信息
    _, err = tx.Exec("INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (?, ?, ?)", fromID, toID, amount)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("记录交易失败: %v", err)
    }

    // 提交事务
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("提交事务失败: %v", err)
    }

    return nil
}

func main() {
    // 数据库连接配置
    dsn := "root:123456@tcp(127.0.0.1:3306)/dbtest?charset=utf8mb4&parseTime=True"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("数据库连接失败: ", err)
    }
    defer db.Close()

    // 测试数据库连接
    err = db.Ping()
    if err != nil {
        log.Fatal("数据库 Ping 失败: ", err)
    }

    // 设置连接池参数（可选但推荐）
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)

    // 执行转账：从账户1转账100元到账户2
    err = transferMoney(db, 1, 2, 100.0)
    if err != nil {
        log.Fatal("转账失败: ", err)
    }

    fmt.Println("转账成功！")
}