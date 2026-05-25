// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface ITreasuryNative {
    error InsufficientBalance(uint256 requested, uint256 available);
    error ZeroAmount();
    error TransferFailed(address to, uint256 amount);

    event NativeDeposited(
        address indexed from,
        uint256 amount,
        uint256 balanceAfter
    );

    event NativeWithdrawn(
        address indexed to,
        uint256 amount,
        uint256 balanceAfter
    );

    function depositNative() external payable;
    function withdrawNative(uint256 amount) external;
    function nativeBalance() external view returns (uint256);
}