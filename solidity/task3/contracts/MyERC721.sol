// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyERC721 is Ownable, ERC721 {

    constructor(
        string memory name,
        string memory symbol
    ) Ownable(msg.sender) ERC721(name, symbol) {}

    function mint(uint256 tokenId) external {
        super._mint(msg.sender, tokenId);
    }
}
