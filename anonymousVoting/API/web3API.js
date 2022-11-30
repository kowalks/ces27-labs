const Web3 = require('web3')
const solc = require('solc')
const fs = require('fs')

// Connect to Ganache
const web3 = new Web3(new Web3.providers.HttpProvider('HTTP://127.0.0.1:8545'))

// Get contract file
const contractFile = fs.readFileSync('../contracts/Election.sol', 'utf-8')

// Input structure for solidity compiler
const input = {
    language: "Solidity",
    sources: {
        "Election.sol": {
            content: contractFile,  
        },
    },
    settings: {
        outputSelection: {
            "*": {
                "*": ["*"],
            },
        },
    },
};

// Compile contract
const output = JSON.parse(solc.compile(JSON.stringify(input)));

// Getting ABI and bytecode
const ABI = output.contracts["Election.sol"]["Election"].abi;
const bytecode = output.contracts["Election.sol"]["Election"].evm.bytecode.object;

// Create the contract
const contract = new web3.eth.Contract(ABI);

// Some account address
const address = "0xA0305EE8C7ECD8F62Aacba32Ce86C4483d05b5A2"

contract.deploy({ data: bytecode })
	.send({ from: address, gas: 1000000 })
	.on("receipt", (receipt) => {
		console.log("Contract Address:", receipt.contractAddress);
	})
	.then((initialContract) => {
		initialContract.methods.vote(1, true).call((err, data) => {
			if (err) console.log(err);
			console.log(data)
		});
	});
