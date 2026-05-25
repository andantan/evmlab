// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

abstract contract TreasuryContext {
    mapping(address => mapping(address => uint256)) internal allocation;
    mapping(address => mapping(address => uint256)) internal pending;
    mapping(address => uint256) internal totalPending;
    mapping(address => uint256) internal totalAllocation;

    function _available(address asset, uint256 balance) internal view returns (uint256) {
        uint256 reserved = totalPending[asset] + totalAllocation[asset];

        if (balance <= reserved) {
            return 0;
        }

        unchecked {
            return balance - reserved;
        }
    }
}
