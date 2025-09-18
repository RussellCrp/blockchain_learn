// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

contract MyNFT is IERC721{

    string public _name;
    string public _symbol;

    mapping(address owner => uint256) private _balances;

    mapping(uint256 tokenId => address ) private _owners;

    mapping(uint256 tokenId => address ) private _tokenApprovals;

    mapping(address owner => mapping(address operater => bool)) private _operatorApprovals;

    constructor(string memory name_, string memory symbol_){
        _name = name_;
        _symbol = symbol_;
    }

    modifier check0Address(address addr){
        require(addr != address(0), "address invalid");
        _;
    }

    function balanceOf(address owner) external view returns (uint256 balance) {
        return _balances[owner];
    }

    function ownerOf(uint256 tokenId) external view returns (address owner) {
        return _owners[tokenId];
    }

    function safeTransferFrom(address from, address to, uint256 tokenId, bytes calldata data) external {
        transferFrom(from, to, tokenId);
    }

    function safeTransferFrom(address from, address to, uint256 tokenId) external {
        transferFrom(from, to, tokenId);
    }

    function transferFrom(address from, address to, uint256 tokenId) public check0Address(from) check0Address(to) {
        address own = _owners[tokenId];
        require(own != address(0) && own == from, "tokenId invalid");

        if (msg.sender != from) {
            address approve_ = _tokenApprovals[tokenId];
            bool isAllApprovals = _operatorApprovals[from][msg.sender];
            require(msg.sender == approve_ || isAllApprovals, string.concat("no auth, tokenId:", string(abi.encodePacked(tokenId))));
        }

        _balances[from]--;
        _owners[tokenId] = to;
        _tokenApprovals[tokenId] = address(0);
        emit Transfer(from, to, tokenId);
    }

    function approve(address to, uint256 tokenId) external check0Address(to){
        address own = _owners[tokenId];
        require(own == msg.sender, "tokenId invalid");
        _tokenApprovals[tokenId] = to;
        emit Approval(msg.sender, to, tokenId);
    }

    function setApprovalForAll(address operator, bool approved) external check0Address(operator){
        _operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function getApproved(uint256 tokenId) external view returns (address operator) {
        address owner = _owners[tokenId];
        require(owner != address(0), "tokenId invalid");
        address operater = _tokenApprovals[tokenId];
        require(operater != address(0), "tokenId no operator");
        return operater;
    }

    function isApprovedForAll(address owner, address operator) external view check0Address(owner) check0Address(operator) returns (bool) {
        return _operatorApprovals[owner][operator];
    }

    function mint(address to, uint256 tokenId) external {
        address owner = _owners[tokenId];
        require(owner == address(0), "tokenId exists");
        _owners[tokenId] = to;
        _balances[to]++;
    }

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return false;
    }
}
