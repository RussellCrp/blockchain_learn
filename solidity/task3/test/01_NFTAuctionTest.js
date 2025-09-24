
const { ethers, deployments, upgrades } = require("hardhat");
const { expect } = require("chai");

describe("Test Auction", async () => {

    const tokenID = 999
    const auctionID = 1


    let erc721Contract, auctionFactoryContractProxy, auctionContractProxy
    let buyer1, buyer2, seller, admin


    it("deploy", async () => {
        [buyer1, seller, admin, buyer2] = await ethers.getSigners()
        console.log("admin::", admin.address)
        console.log("seller::", seller.address)
        console.log("buyer1::", buyer1.address)
        console.log("buyer2::", buyer2.address)
        // 获取 erc721 合约的工厂实例
        const nft = await ethers.getContractFactory("MyERC721")
        // 指定用 admin 部署合约（而非默认的 account0）
        erc721Contract = await nft.connect(admin).deploy("name", "symbol")
        // 等待合约部署上链
        await erc721Contract.waitForDeployment()
        // 获取已部署合约的地址
        let erc721ContractAddr = await erc721Contract.getAddress()
        expect(erc721ContractAddr).to.length.greaterThan(0)

        // 获取拍卖合约的工厂实例
        const auctionFactory = await ethers.getContractFactory("AuctionFactory")
        // 使用 透明代理模式部署拍卖合约，传入构造函数参数
        auctionFactoryContractProxy = await upgrades.deployProxy(auctionFactory, [admin.address],
            { initializer: "initialize", kind: "transparent", signer: admin })
        
        // 等待部署上链
        await auctionFactoryContractProxy.waitForDeployment()
        expect(await auctionFactoryContractProxy.admin()).to.equal(admin.address)
        const auctionFactoryProxyAddr = await auctionFactoryContractProxy.getAddress()

        console.log("erc721合约地址::", erc721ContractAddr)
        console.log("工厂合约代理地址::", auctionFactoryProxyAddr)
    })

    it("mint nft", async ()=>{
        // 铸造nft
        await erc721Contract.connect(seller).mint(tokenID)
        // 获取 nft 的拥有者
        const nftOwner = await erc721Contract.ownerOf(tokenID)
        expect(nftOwner).to.equal(seller.address)
    })

    it("create auction", async ()=>{
        // 创建拍卖合约
        let erc721ContractAddr = await erc721Contract.getAddress()
        const tx = await auctionFactoryContractProxy.connect(admin)
            .createAuction(auctionID, seller.address, tokenID, 100, 60 * 5, erc721ContractAddr)
        await tx.wait()
        
        // 通过地址重新绑定合约实例
        const auctionProxyAddress = await auctionFactoryContractProxy.auctions(auctionID)
        auctionContractProxy = await ethers.getContractAt("NFTAuction", auctionProxyAddress);
        
        
        const auctionTokenID = await auctionContractProxy.tokenID()
        expect(auctionTokenID).to.equal(tokenID)
        
        console.log("拍卖合约代理地址::", auctionProxyAddress)
    })

    it("upgrade auction", async ()=>{
        // 将代理合约注册到 Upgrades 插件
        const LogicContract = await ethers.getContractFactory("NFTAuction");
        await upgrades.forceImport(auctionContractProxy, LogicContract, {
            kind: "transparent", // 或 "uups"（根据代理类型）
        });
        // 获取拍卖合约工厂实例
        const auctionV2 = await ethers.getContractFactory("NFTAuctionV2")
        // 使用admin升级拍卖合约
        let auctionProxyAdress = await auctionContractProxy.getAddress()
        const oldAddress = await upgrades.erc1967.getImplementationAddress(auctionProxyAdress)
        auctionContractProxy = await upgrades.upgradeProxy(auctionProxyAdress, auctionV2.connect(admin))
        await auctionContractProxy.waitForDeployment()
        // 验证
        const newAddress = await upgrades.erc1967.getImplementationAddress(auctionProxyAdress)
        expect(oldAddress).to.not.equal(newAddress)
        expect(await auctionContractProxy.tokenID()).to.equal(tokenID)
        expect(await auctionContractProxy.version()).to.equal("v2")
    })


})