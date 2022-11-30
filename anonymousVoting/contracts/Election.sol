// SPDX-License-Identifier: MIT

pragma solidity ^0.8.17;

contract Election {

    // Voter model
    struct Voter {
        uint id;
        string name;
        bool voted;
    }

    // Vote model
    struct Vote {
        address voterAddress;
        bool choice;
    }

    // Store accounts that have voted
    mapping(uint => Voter) public voters;

    // Store votes
    mapping(address => Vote) private votes;

    // Store votes count
    uint public finalizedVotes = 0;

    // Auxiliary count
    uint public aux = 0;

    // Store voters amount
    uint public votersAmount;

   constructor() {
        addCandidate("Einstein");
        addCandidate("Newton");
        addCandidate("Tesla");
        votersAmount = 3;
    }

    function addCandidate (string memory _name) private {
        aux++;
        voters[aux] = Voter(aux, _name, false);
    }

    function vote (uint _candidateId, bool choice) public {
        // require that they haven't voted before
        require(!voters[_candidateId].voted);

        // require a valid candidate
        require(_candidateId > 0 && _candidateId <= votersAmount);

        // record that voter has voted
        voters[_candidateId].voted = true;

        // record vote
        votes[msg.sender] = Vote(msg.sender, choice);

        // update candidate vote Count
        finalizedVotes++;
    }
} 