
// deploy/00_deploy_my_contract.js
// export a function that get passed the Hardhat runtime environment
const { deployments, upgrades, ethers } = require("hardhat")
require("dotenv").config();

const fs = require("fs");
const path = require("path");

module.exports = async ({ getNamedAccounts, deployments }) => {
    const erc721Contract = await deployERC721();
    const erc721Address = await erc721Contract.getAddress()
    console.log("erc721合约地址", erc721Address)

    const auctionFactoryProxy = await deployNFTFactoryProxy();
    const proxyAddress = await auctionFactoryProxy.getAddress()
    console.log("auctionFactory代理合约地址", proxyAddress)
    
    
    
    // 保存部署信息到文件，供前端调用
    const { save } = deployments;
    await save("depolyInfo", {
        erc721Address,
        proxyAddress,
        abi: {
            erc721: erc721Contract.interface.format("json"),
            auctionFactory: auctionFactoryProxy.interface.format("json"),
        }
    })
    console.log("署信息已保存到 deployments/depolyInfo.json")
};

async function deployERC721() {
    const erc721 = await ethers.getContractFactory("MyERC721")
    // 指定用 admin 部署合约（而非默认的 account0）
    const erc721Contract = await erc721.deploy("MyERC721", "MEC")
    // 等待合约部署上链
    await erc721Contract.waitForDeployment()
    // 获取已部署合约的地址
    return erc721Contract
}

async function deployNFTFactoryProxy() {
    const auctionFactory = await ethers.getContractFactory("AuctionFactory")
    const proxy = await upgrades.deployProxy(auctionFactory, [process.env.pbk], { initializer: "initialize", })
    // 等待合约部署上链
    await proxy.waitForDeployment()
    return proxy
}

// add tags and dependencies
module.exports.tags = ["depoly"];