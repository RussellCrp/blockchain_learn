// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "hardhat/console.sol"; // 导入 Hardhat 控制台日志库

contract NFTAuctionV2 is Initializable {
    uint256 public auctionID;
    // 卖家
    address public seller;
    // NFT合约地址
    address public erc721Addr;
    // NFT TokenID
    uint256 public tokenID;
    // 起拍价
    uint256 public startPrice;
    // 结束时间
    uint256 public endTime;
    // 最高出价
    uint256 public highestBid;
    // 最高出价者
    address public highestBidder;
    // 最高出价 erc20 地址, address(0) 代表eth
    address public highestERC20Addr;


    // chainlink 预言机地址
    mapping(address erc20 => AggregatorV3Interface priceFeed) public priceFeeds;

    // erc20地址 未喂价
    error ERC20FeedUndefined(address erc20);


    // 初始合约
    function initialize(bytes memory params) external initializer {
        (address _seller,uint256 _tokenID, uint256 _startPrice,uint256 duration, address _erc721Addr, uint256 _auctionID) = abi.decode(params, (address, uint256, uint256, uint256, address, uint256));
        initAuction(_seller, _tokenID, _startPrice, duration, _erc721Addr, _auctionID);
    }

    // 初始化拍卖
    function initAuction(address _seller, uint256 _tokenID, uint256 _startPrice, uint256 duration, address _erc721Addr, uint256 _auctionID) internal {
        seller = _seller;
        tokenID = _tokenID;
        startPrice = _startPrice;
        endTime = duration + block.timestamp;
        erc721Addr = _erc721Addr;
        auctionID = _auctionID;
    }

    // 出价
    function bidPrice(address erc20Addr_, uint256 amount_) external payable {
        require(block.timestamp < endTime, "EOA");

        // 计算出价价值
        uint256 payableAmount;
        if (erc20Addr_ != address(0)) {
            payableAmount = amount_ * uint(getChainlinkDataFeedLatestAnswer(erc20Addr_));
        } else {
            payableAmount = msg.value * uint(getChainlinkDataFeedLatestAnswer(address(0)));
        }

        // 计算当前价值
        uint256 currentBidUSD;
        if (highestBid != 0){
            currentBidUSD = highestBid * uint(getChainlinkDataFeedLatestAnswer(highestERC20Addr));
        }
        require(payableAmount > currentBidUSD && payableAmount > startPrice, "The bid is lower than the highest bid price");

        address highestBidderBefore = highestBidder;
        address highestERC20AddrBefore = highestERC20Addr;
        uint256 highestBidBefore = highestBid;

        // 修改最高出价信息
        highestBidder = msg.sender;
        highestERC20Addr = erc20Addr_;
        if (erc20Addr_ == address(0)) {
            highestBid = msg.value;
        } else {
            highestBid = amount_;
            bool ok = IERC20(erc20Addr_).transferFrom(msg.sender, address(this), amount_);
            if (!ok) {
                revert("erc20 token transfer failed");
            }
        }

        // 退还给之前最高出价者
        if (highestERC20AddrBefore == address(0)) {
            payable(highestBidderBefore).transfer(highestBidBefore);
        } else {
            bool ok = IERC20(highestERC20AddrBefore).transfer(highestBidderBefore, highestBidBefore);
            if (!ok) {
                revert("The refund to the previous highest bidder fails");
            }
        }

    }

    // 结束拍卖
    function endAuction() external payable {
        require(block.timestamp >= endTime, "The auction is not over.");
        IERC721(erc721Addr).safeTransferFrom(seller, highestBidder, tokenID);
        if (highestERC20Addr != address(0)) {
            bool ok = IERC20(highestERC20Addr).transfer(seller, highestBid);
            if (!ok) {
                revert("seller erc20 failed collection");
            }
        } else {
            payable(seller).transfer(highestBid);
        }
    }

    // ETH / USD ==>    0x694AA1769357215DE4FAC081bf1f309aDC325306
    // USDC / USD ===>  0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E
    // 喂价
    function setPriceFeed(address tokenAddress, address _priceFeed) public {
        priceFeeds[tokenAddress] = AggregatorV3Interface(_priceFeed);
    }

    // 获取最新币对价格
    function getChainlinkDataFeedLatestAnswer(address tokenAddress) public view returns (int) {
        AggregatorV3Interface priceFeed = priceFeeds[tokenAddress];
        if (address(priceFeed) == address(0)){
            revert ERC20FeedUndefined(tokenAddress);
        }
        // prettier-ignore
        (
        /* uint80 roundId */,
            int256 answer,
        /*uint256 startedAt*/,
        /*uint256 updatedAt*/,
        /*uint80 answeredInRound*/
        ) = priceFeed.latestRoundData();
        return answer;
    }

    function version() external pure returns (string memory) {
        return "v2";
    }
}