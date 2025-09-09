-- 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和
-- transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
-- 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
-- 如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
DO $$
BEGIN
-- 检查余额
IF (SELECT balance FROM accounts WHERE id = 1) >= 100 THEN
    -- 扣款
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;
    -- 入账
    UPDATE accounts SET balance = balance + 100 WHERE id = 2;
    -- 记录转账
    INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (1, 2, 100);
    COMMIT;
    RAISE NOTICE '转账成功';
ELSE
    ROLLBACK;
    RAISE NOTICE '转账失败：余额不足';
END IF;
END $$;