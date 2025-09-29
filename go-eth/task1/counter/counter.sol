// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract counter {
    uint256 public num;

    function count() external returns (uint256) {
        num++;
        return num;
    }
}
