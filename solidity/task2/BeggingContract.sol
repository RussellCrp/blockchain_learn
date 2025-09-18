// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract BeggingContract {

    address[3] private _topThree;
    uint256 private _totalETH;
    address private _cyberbeggar;
    mapping(address account => uint256) private _donor;
    event Donation(address from, uint256 amount);

    uint256 public endTime;

    constructor(uint256 duringTime){
        _cyberbeggar = msg.sender;
        endTime = duringTime + block.timestamp;
    }

    modifier onlyOwner() {
        require(msg.sender == _cyberbeggar, "not owner");
        _;
    }

    modifier checkEnd(){
        require(block.timestamp < endTime, "beg end");
        _;
    }

    function donate() external payable checkEnd {
        require(msg.value > 0, "Must send ETH");
        _donor[msg.sender] += msg.value;
        _totalETH += msg.value;

        topThreeSort();

        emit Donation(msg.sender, msg.value);
    }

    function topThreeSort() internal virtual {
        if (msg.value <= _donor[_topThree[_topThree.length - 1]]){
            return;
        }
        uint256 tmpVal = msg.value;
        address tmpAddr = msg.sender;
        for(uint256 i = 0; i < _topThree.length; i++){
            address topAddr = _topThree[i];
            uint256 topAmount = _donor[topAddr];
            if (tmpVal > topAmount){
                _topThree[i] = tmpAddr;
                tmpVal = topAmount;
                tmpAddr = topAddr;
            }
        }
    }


    function withdraw(address payable _to) external onlyOwner{
        _to.transfer(_totalETH); // 发送 ETH 给 _to
        _totalETH = 0;
    }

    function getDonation(address donorAddr) public view returns (uint256) {
        return _donor[donorAddr]  / 1 ether;
    }

    function getTopThree() public view returns (address[3] memory) {
        return _topThree;
    }
}
