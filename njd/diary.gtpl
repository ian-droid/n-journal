<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="utf-8">
  <title>Diary</title>
  <script></script>

</head>
<body>
  <h2>Recent Activities</h2>
  <table>
    <tr><th>Date</th><th>Content</th></tr>
    {{range .Diary}}
      {{if .Highlighted}}
        <tr class="highlighted">
      {{else}}
        <tr>
      {{end}}
        <td>{{.Oid}}</td><td>{{.Date}}</td><td>{{.Content}}</td></tr>
    {{end}}
  </table>
  <h2>New Diary</h2>
  <form action="diary" method="POST">
    <label for="date">Date:</label>
    <input type="date" name="date" />
    <label for="highlighted">Highlighted:</label>
    <input type="checkbox" name="highlighted" /><br />
    <label for="content">Content:</label>
    <textarea name="content"></textarea><br />
    <input type="submit"></input>
  </form>
</body>
</html>
