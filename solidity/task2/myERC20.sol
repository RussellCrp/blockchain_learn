// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
    合约包含以下标准 ERC20 功能：
    balanceOf：查询账户余额。
    transfer：转账。
    approve 和 transferFrom：授权和代扣转账。
    使用 event 记录转账和授权操作。
    提供 mint 函数，允许合约所有者增发代币。
    提示：
    使用 mapping 存储账户余额和授权信息。
    使用 event 定义 Transfer 和 Approval 事件。
    部署到sepolia 测试网，导入到自己的钱包
*/
contract myERC20 {

    uint256 public totalSupply;
    address private owner;
    mapping(address => uint256) private balance;
    mapping(address => mapping(address => uint256)) private _allowances;

    event Transfer(address indexed from, address indexed to, uint256 value);

    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor(){
        owner = msg.sender;
    }

    function balanceOf(address account) external view returns(uint256){
        return balance[account];
    }

    function transfer(address to, uint256 value) external returns (bool) {
        require(address(0) != to, "addr invalid");
        require(balance[msg.sender] >= value, "balance insufficient");
        balance[msg.sender] -= value;
        balance[to] += value;
        emit Transfer(msg.sender, to, value);
        return true;
    }

    function approve(address spender, uint256 value) external returns (bool){
        require(address(0) != spender, "addr invalid");
        require(balance[msg.sender] >= value, "balance insufficient");
        _allowances[msg.sender][spender] += value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    function transferFrom(address from, address to, uint256 value) external returns (bool){
        require(address(0) != from, "from addr invalid");
        require(address(0) != to, "to addr invalid");
        require(balance[from] >= value, "from balance insufficient");
        require(_allowances[from][msg.sender] >= value, "approve balance insufficient");
        balance[from] -= value;
        _allowances[from][msg.sender] -= value;
        balance[to] += value;
        emit Transfer(from, to, value);
        return true;
    }

    function mint(address account, uint value) external {
        require(msg.sender == owner, "not allowed");
        balance[account] += value;
        totalSupply += value;
    }

}