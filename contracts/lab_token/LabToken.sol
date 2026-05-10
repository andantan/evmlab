// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "./interfaces/IERC20.sol";

contract LabToken is IERC20 {
    string public name;
    string public symbol;
    uint8 public decimals;
    
    uint256 private _totalSupply;

    mapping(address => uint256) private balances;
    mapping(address => mapping(address => uint256)) private allowances;

    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decmials,
        uint256 initialSupply
    ){
        name = _name;
        symbol = _symbol;
        decimals = _decmials;

        _mint(msg.sender, initialSupply);
    }

    function totalSupply() external view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) external view returns (uint256) {
        if (account == address(0)) {
            revert InvalidAddress();
        }

        return balances[account];
    }

    function transfer(address to, uint256 value) external returns (bool) {
        if (to == address(0)) {
            revert InvalidAddress();
        }

        _transfer(msg.sender, to, value);
        return true;
    }

    function allowance(address owner, address spender) external view returns (uint256) {
        if (owner == address(0) || spender == address(0)) {
            revert InvalidAddress();
        }

        return allowances[owner][spender];
    }

    function approve(address spender, uint256 value) external returns (bool) {
        if (spender == address(0)) {
            revert InvalidAddress();
        }

        allowances[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    function transferFrom(address from, address to, uint256 value) external returns (bool) {
        uint256 currentAllowance = allowances[from][msg.sender];

        if (currentAllowance < value) {
            revert InsufficientAllowance(from, msg.sender, value, currentAllowance);
        }

        unchecked {
            allowances[from][msg.sender] = currentAllowance - value;
        }

        emit Approval(from, msg.sender, allowances[from][msg.sender]);
        _transfer(from, to, value);
        return true;

    }

    function _transfer(address from, address to, uint256 value) internal {
        if (from == address(0) || to == address(0)) {
            revert InvalidAddress();
        }

        uint256 fromBalance = balances[from];

        if (fromBalance < value) {
            revert InsufficientBalance(from, value, fromBalance);
        }

        unchecked {
            balances[from] = fromBalance - value;
        }

        balances[to] += value;
        emit Transfer(from, to, value);
    }

    function _mint(address to, uint256 value) internal {
        if (to == address(0)) {
            revert InvalidAddress();
        }

        _totalSupply += value;
        balances[to] += value;

        emit Transfer(address(0), to, value);
    }
}
