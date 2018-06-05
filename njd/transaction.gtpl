<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="utf-8">
  <title>Transactions</title>
  <link rel="stylesheet" href="style.css">
  <script></script>
</head>
<body>
  <h2>Recent Transactions</h2>
  <table>
    <tr>
      <th>Date</th><th>Item</th><th>Description</th><th>Direction</th>
      <th>Currency</th><th>Amount</th><th>Payment</th><th>Bank</th>
    </tr>
    {{range .Transaction}}
      <tr>
        <td>{{.Date}}</td>
        <td>{{.Item}}</td>
        <td>{{.Description}}</td>
        <td>{{.Direction}}</td>
        <td>{{.CurrencyName}}</td>
        <td>{{.Amount}}</td>
        <td>{{.PaymentName}}</td>
        <td>{{.BankName}}</td>
      </tr>
    {{end}}
  </table>
  <h2>New Transaction</h2>
  <form action="transaction" method="POST">
    <label for="date">Date:</label>
    <input type="date" name="date" />
    <label for="item">Item:</label>
    <input type="text" name="item" />
    <label for="description">Description:</label>
    <input type="text" name="description" /><br />
    <!--<label for="currency">Currency:</label>-->
    <select name="currency">
      {{range .Currency}}
      <option value="{{.Id}}">{{.Prefix}}</option>
      {{end}}
    </select>
    <!--<label for="amount">Amount:</label>-->
    <input type="text" name="amount" />
    <input type="radio" name="direction" value="pay" checked>Pay</input>
    <input type="radio" name="direction" value="income">Income</input>
    <label for="payment">Payment:</label>
    <select name="payment">
      {{range .Payment}}
      <option value="{{.Id}}" {{if .Priority}}selected{{end}}>{{.Name}}</option>
      {{end}}
    </select>
    <label for="bank">Bank:</label>
    <select name="bank">
      {{range .Bank}}
      <option value="{{.Id}}" {{if .Priority}}selected{{end}}>{{.Name}}</option>
      {{end}}
    </select>
    <input type="submit" value="{{.Mode}}" />
  </form>
</body>
</html>
