// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library Asset {
    address internal constant NATIVE = address(0);

    function isNative(address asset) internal pure returns (bool) {
        return asset == NATIVE;
    }

    function isToken(address asset) internal pure returns (bool) {
        return asset != NATIVE;
    }
}