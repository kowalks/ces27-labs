// SPDX-License-Identifier: MIT

pragma solidity >=0.8;

import "anonymousVoting/contracts/EllipticCurve.sol";

contract AnonymousVoting {
    // EllipticCurve y^2 = x^3 + a x^2 + b
    // Over Galois Field GF(p)
    // Name: secp256k1
    uint256 public constant X = 12;
    uint256 public constant Y = 15;


    function publicKey(uint privKey) external pure returns (uint, uint) {
        return EllipticCurve.pMul(X, Y, privKey);
    }
}
