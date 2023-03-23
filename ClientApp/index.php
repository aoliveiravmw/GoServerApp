<!DOCTYPE html>
<html>
<head>
	<title>Add Plant</title>
</head>
<body>
	<h1>Add Plant</h1>
	<form action="/addplant" method="POST">
		<label for="plant_name">Plant Name:</label>
		<input type="text" name="plant_name"><br>

		<label for="watering_interval_hours">Watering Interval (in hours):</label>
		<input type="number" name="watering_interval_hours"><br>

		<label for="last_watered">Last Watered:</label>
		<input type="date" name="last_watered"><br>

		<input type="submit" value="Add Plant">
	</form>
</body>
</html>
