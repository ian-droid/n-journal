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
        <td>{{.Date}}</td><td>{{.Content}}</td><td><a href="diary?Edit={{.Oid}}">Edit</a></td></tr>
    {{end}}
  </table>
  {{if .Message}}<div id="Message">{{.Message}}</div>{{end}}
  <h2>{{.Mode}} Diary</h2>
  {{if .Diary2Update.Oid}}
  <form action="diary" method="POST">
    <label for="oid">Oid:</label>
    <input type="text" name="oid" value="{{.Diary2Update.Oid}}" readonly/>
    <label for="date">Date:</label>
    <input type="date" name="date" value="{{.Diary2Update.Date}}" readonly/>
    <label for="highlighted">Highlighted:</label>
    <input type="checkbox" name="highlighted" {{if .Diary2Update.Highlighted}}checked{{end}}/><br />
    <label for="content">Content:</label><br />
    <textarea name="content">{{.Diary2Update.Content}}</textarea><br />
    <input type="submit" />
  </form>
  {{else}}
  <form action="diary" method="POST">
    <label for="date">Date:</label>
    <input type="date" name="date" />
    <label for="highlighted">Highlighted:</label>
    <input type="checkbox" name="highlighted" /><br />
    <label for="content">Content:</label><br />
    <textarea name="content"></textarea><br />
    <input type="submit" />
  </form>
  {{end}}
</body>
</html>
