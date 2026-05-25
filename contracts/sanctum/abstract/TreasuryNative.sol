// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {NexusContext} from "./NexusContext.sol";
import {TreasuryContext} from "./TreasuryContext.sol";
import {ITreasuryNative} from "../interfaces/ITreasuryNative.sol";
import {Asset} from "../libraries/Asset.sol";

abstract contract TreasuryNative is NexusContext, TreasuryContext, ITreasuryNative {
    receive() external payable onlyMember {
        _depositNative(msg.sender, msg.value);
    }

    function depositNative() external payable onlyMember {
        _depositNative(msg.sender, msg.value);
    }

    function requestNative(uint256 amount) external onlyMember {
        _requestNative(msg.sender, amount);
    }

    function approveNative(address user, uint256 amount) external onlyMaster {
        _approveNative(user, amount);
    }

    function approveNativeAll(address user) external onlyMaster {
        _approveNative(user, pending[user][Asset.NATIVE]);
    }

    function withdrawNative(uint256 amount) external onlyMember {
        _withdrawNative(msg.sender, amount);
    }

    function withdrawNativeAll() external onlyMember {
        _withdrawNative(msg.sender, allocation[msg.sender][Asset.NATIVE]);
    }

    function nativeBalance() external view returns (uint256) {
        return address(this).balance;
    }

    function nativeAvailable() external view returns (uint256) {
        return _available(Asset.NATIVE, address(this).balance);
    }

    function nativeAllocation(address user) external view returns (uint256) {
        return allocation[user][Asset.NATIVE];
    }

    function nativePending(address user) external view returns (uint256) {
        return pending[user][Asset.NATIVE];
    }

    function _depositNative(address from, uint256 amount) internal {
        if (amount == 0) {
            revert ZeroAmount();
        }

        emit Deposited(from, Asset.NATIVE, amount, address(this).balance);
    }

    function _requestNative(address from, uint256 amount) internal {
        if (amount == 0) {
            revert ZeroAmount();
        }

        address native = Asset.NATIVE;

        uint256 available = _available(native, address(this).balance);
        if (available < amount) {
            revert InsufficientBalance(native, amount, available);
        }

        uint256 userPending = pending[from][native] + amount;
        pending[from][native] = userPending;
        totalPending[native] += amount;

        emit Requested(from, native, amount, userPending);
    }

    function _approveNative(address user, uint256 amount) internal {
        if (amount == 0) {
            revert ZeroAmount();
        }

        address native = Asset.NATIVE;

        uint256 userPending = pending[user][native];
        if (userPending == 0) {
            revert NoPendedBalance(native, user);
        }
        if (userPending < amount) {
            revert InsufficientPending(native, amount, userPending);
        }

        uint256 userAllocation = allocation[user][native] + amount;
        unchecked {
            pending[user][native] = userPending - amount;
            totalPending[native] -= amount;
        }
        allocation[user][native] = userAllocation;
        totalAllocation[native] += amount;

        emit Approved(user, native, amount, userAllocation);
    }

    function _withdrawNative(address to, uint256 amount) internal {
        if (amount == 0) {
            revert ZeroAmount();
        }

        address native = Asset.NATIVE;

        uint256 allocated = allocation[to][native];
        if (allocated == 0) {
            revert NoAllocatedBalance(native, to);
        }
        if (allocated < amount) {
            revert InsufficientAllocation(native, amount, allocated);
        }

        uint256 balance = address(this).balance;
        if (balance < amount) {
            revert InsufficientBalance(native, amount, balance);
        }

        uint256 balanceAfter;
        unchecked {
            balanceAfter = balance - amount;
            allocation[to][native] = allocated - amount;
            totalAllocation[native] -= amount;
        }

        (bool ok,) = to.call{value: amount}("");
        if (!ok) {
            revert TransferFailed(native, to, amount);
        }

        emit Withdrawn(to, native, amount, balanceAfter);
    }
}
