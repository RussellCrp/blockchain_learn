// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

// import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./NFTAuction.sol";
import "hardhat/console.sol"; // 导入 Hardhat 控制台日志库
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";


contract AuctionFactory  is Initializable{
    // 每一个拍卖对应一个合约
    mapping(uint256 auctionID => address auction) public auctions;
    address public admin;

    modifier onlyOwner() {
        require(admin == msg.sender, "Ownable: caller is not the owner");
        _;
    }

    function initialize(address _admin) external initializer {
        admin = _admin;
    }

    function createAuction(uint256 auctionID, address seller, uint256 tokenID, uint256 startPrice, uint256 duration, address erc721Addr) external onlyOwner returns (address) {
        // 用create部署新合约
        bytes memory params = abi.encode(seller, tokenID, startPrice, duration, erc721Addr, auctionID);
        
        NFTAuction auction = new NFTAuction();
        TransparentUpgradeableProxy auctionProxy = new TransparentUpgradeableProxy(
                address(auction),
                admin,
                abi.encodeWithSignature("initialize(bytes)", params)
            );
        address proxy = address(auctionProxy);
        // 记录拍卖合约地址
        auctions[auctionID] = proxy;
        return proxy;
    }

    
}
