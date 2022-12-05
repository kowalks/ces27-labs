// SPDX-License-Identifier: MIT

pragma solidity >=0.8;

library EllipticCurve {
    // EllipticCurve y^2 = x^3 + a x^2 + b
    // Over Galois Field GF(p)
    // Name: secp256k1
    uint256 public constant A = 70;
    uint256 public constant B = 0;
    uint256 public constant P = 71;

    // invMod by Fermat Little Theorem: a^(p-2) is the inverse of a mod p
    function modInv(uint256 a) internal pure returns (uint256) {
        return modExp(a, P - 2);
    }

    // computes a^b mod p in O(log b)
    function modExp(uint256 a, uint256 b) internal pure returns (uint256) {
        if (b == 0) return 1;

        uint256 c = modExp(a, b / 2);
        c = (c * c) % P;
        if (b % 2 != 0) c = (c * a) % P;
        return c;
    }

    function modMul(uint256 a, uint256 b) internal pure returns (uint256) {
        return (a * b) % P;
    }

    function modDiv(uint256 a, uint256 b) internal pure returns (uint256) {
        return modMul(a, modInv(b));
    }

    function modAdd(uint256 a, uint256 b) internal pure returns (uint256) {
        return (a + b) % P;
    }

    function modSub(uint256 a, uint256 b) internal pure returns (uint256) {
        return (a + P - b) % P;
    }

    function pAdd(
        uint256 xp,
        uint256 yp,
        uint256 xq,
        uint256 yq
    ) internal pure returns (uint256, uint256) {
        if (xp == xq) {
            return pDbl(xp, yp);
        }

        uint256 lambda = modDiv(modSub(yq,yp), modSub(xq,xp));
        uint256 xr = modSub(modExp(lambda, 2), modAdd(xp, xq));
        uint256 yr = modSub(modMul(lambda, modSub(xp, xr)),yp);
        return (xr, yr);
    }

    function pDbl(uint256 xp, uint256 yp)
        internal
        pure
        returns (uint256, uint256)
    {
        uint256 xq = xp;

        uint256 lambda = modDiv(3 * modExp(xp, 2) + A, modMul(yp, 2));
        uint256 xr = modSub(modExp(lambda, 2), modAdd(xp, xq));
        uint256 yr = modMul(lambda, modSub(xp, xr)) - yp;
        return (xr, yr);
    }

    function pNeg(uint256 x, uint256 y)
        internal
        pure
        returns (uint256, uint256)
    {
        return (x, modSub(0, y));
    }

    function pMul(
        uint256 x,
        uint256 y,
        uint256 k
    ) internal pure returns (uint256, uint256) {
        if (k == 1) return (x, y);

        uint256 xx;
        uint256 yy;
        (xx, yy) = pMul(x, y, k / 2);
        (xx, yy) = pDbl(xx, yy);

        if (k % 2 != 0) (xx, yy) = pAdd(xx, yy, x, y);
        return (xx, yy);
    }
}
