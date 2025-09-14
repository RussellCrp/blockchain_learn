// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Voting {

    mapping(address => uint) private voteMapping;
    mapping(address => bool) private hasVoted;
    address[] private candidates;

    function vote(address candidate) external {
        require(!hasVoted[msg.sender],  "already voted");
        hasVoted[msg.sender] = true;
        voteMapping[candidate]++;
        if (voteMapping[candidate] == 1){
            candidates.push(candidate);
        }
    }

    function getVotes(address candidate) external view returns (uint) {
        return voteMapping[candidate];
    }

    function resetVotes() external {
        for (uint i = 0; i < candidates.length; i++) {
            address candidate = candidates[i];
            voteMapping[candidate] = 0;
        }
        delete candidates;
    }

}