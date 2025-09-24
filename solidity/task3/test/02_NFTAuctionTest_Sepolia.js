
const { ethers, deployments, upgrades } = require("hardhat");
const { expect } = require("chai");
require("dotenv").config();

describe("Test Auction", async () => {

    
    const ETH_USD = "0x694AA1769357215DE4FAC081bf1f309aDC325306"
    const USDC_USD = "0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E"
    const ERC20_ADDRESS = "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238" // USDC 测试网合约地址

    const erc721ContractAddress = "0xcfAe3eA6EaDd306C6d5aDbA5A6138e1ac06774d4"
    const auctionFactoryContractProxyAddress = "0xEdd036d8cd8B02C9e253CD2e4cB09Df099Afa0B8"
    const auctionContractAddress = "0x07dC4FfCc9bbbFc8053d7872FD458d50e076e8c0"

    const tokenID = 999
    const auctionID = 1
    const sellerAddr = process.env.pbk


    let erc721Contract, auctionFactoryContractProxy, auctionContractProxy
    

    beforeEach(async () => {
        erc721Contract = await ethers.getContractAt("MyERC721", erc721ContractAddress)
        auctionFactoryContractProxy = await ethers.getContractAt("AuctionFactory", auctionFactoryContractProxyAddress)
        auctionContractProxy = await ethers.getContractAt("NFTAuction", auctionContractAddress)
    })

    it("mint nft", async ()=>{
        // 铸造nft
        await erc721Contract.mint(tokenID)
        // 获取 nft 的拥有者
        const nftOwner = await erc721Contract.ownerOf(tokenID)
        expect(nftOwner).to.equal(process.env.pbk)
    })

    it("create auction", async ()=>{
        // 创建拍卖合约 //0x07dC4FfCc9bbbFc8053d7872FD458d50e076e8c0
        const tx = await auctionFactoryContractProxy.createAuction(auctionID, sellerAddr, tokenID, 1, 60 * 5, erc721ContractAddress)
        await tx.wait()
        
        // 通过地址重新绑定合约实例
        const auctionProxyAddress = await auctionFactoryContractProxy.auctions(auctionID)
        auctionContractProxy = await ethers.getContractAt("NFTAuction", auctionProxyAddress);
        
        
        const auctionTokenID = await auctionContractProxy.tokenID()
        expect(auctionTokenID).to.equal(tokenID)
        
        console.log("拍卖合约代理地址::", auctionProxyAddress)
    })
    

    // it("feed price", async ()=>{
    //     // 预设价格
    //     await auctionContractProxy.setPriceFeed("", ETH_USD)
    //     await auctionContractProxy.setPriceFeed(ERC20_ADDRESS, USDC_USD)

    //     expect(await auctionContractProxy.priceFeeds(ERC20_ADDRESS)).to.equal(USDC_USD)
    //     expect(await auctionContractProxy.priceFeeds("")).to.equal(ETH_USD)
    // })

    // it("bid price", async ()=> {
    //     // 竞拍
    //     await auctionContractProxy.bidPrice("", 0,{ value: 3000000 })
    //     expect(await auctionContractProxy.highestBid()).to.equal(eth)

    //     await auctionContractProxy.connect(buyer2).bidPrice("", 0,{ value: 4000000 })
    //     expect(await auctionContractProxy.highestBid()).to.equal(eth)
    // })
})