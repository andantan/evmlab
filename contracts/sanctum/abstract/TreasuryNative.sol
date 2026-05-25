// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SanctumContext} from "./SanctumContext.sol";
import {ITreasuryNative} from "../interfaces/ITreasuryNative.sol";

abstract contract TreasuryNative is SanctumContext, ITreasuryNative {
    receive() external payable onlyMember {
        _depositNative(msg.sender, msg.value);
    }

    function depositNative() external payable onlyMember {
        _depositNative(msg.sender, msg.value);
    }

    function withdrawNative(uint256 amount) external onlyMember {
        _withdrawNative(msg.sender, amount);
    }

    function nativeBalance() external view returns (uint256) {
        return address(this).balance;
    }

    function _depositNative(address from, uint256 amount) internal {
        if (amount == 0) {
            revert ZeroAmount();
        }

        emit NativeDeposited(from, amount, address(this).balance);
    }

    function _withdrawNative(address to, uint256 amount) internal {
        if (amount == 0) {
            revert ZeroAmount();
        }

        uint256 balance = address(this).balance;

        if (balance < amount) {
            revert InsufficientBalance(amount, balance);
        }

        (bool ok,) = to.call{value: amount}("");
        if (!ok) {
            revert TransferFailed(to, amount);
        }

        emit NativeWithdrawn(to, amount, address(this).balance);
    }
}