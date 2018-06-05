<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="utf-8">
  <title>Diary</title>
  <link rel="stylesheet" href="style.css" />
  <script></script>
</head>
<body>
  <h2>Recent Activities</h2>
  <table>
    <tr><th class="date">Date</th><th>Content</th><th>Operation</th></tr>
    {{range .Diary}}
      {{if .Highlighted}}
        <tr class="highlighted">
      {{else}}
        <tr>
      {{end}}
        <td>{{.Date}}</td><td>{{.Content}}</td><td><a href="diary?{{.Oid}}">Edit</a></td></tr>
    {{end}}
  </table>
  <h2>New Diary</h2>
  <form action="diary" method="POST">
    <label for="date">Date:</label>
    <input type="date" name="date" />
    <label for="highlighted">Highlighted:</label>
    <input type="checkbox" name="highlighted" /><br />
    <label for="content">Content:</label><br />
    <textarea name="content"></textarea><br />
    <input type="submit" value="{{.Mode}}"/>
  </form>
</body>
</html>
