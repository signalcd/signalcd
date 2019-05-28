import 'package:angular/angular.dart';

@Pipe('timeago')
class TimeagoPipe extends PipeTransform {
  String transform(DateTime date, [String pattern = "mediumDate"]) {
    if (date == null) {
      return '';
    }

    if (date.isBefore(DateTime.now().subtract(Duration(days: 7)))) {
      return DatePipe().transform(date, pattern);
    }

    if (date.isBefore(DateTime.now().subtract(Duration(days: 1)))) {
      int days = DateTime.now().difference(date).inDays;

      String unit = days > 1 ? 'days' : 'day';
      return '${days} ${unit} ago';
    }

    if (date.isBefore(DateTime.now().subtract(Duration(hours: 1)))) {
      int hours = DateTime.now().difference(date).inHours;

      String unit = hours > 1 ? 'hours' : 'hour';
      return '${hours} ${unit} ago';
    }

    if (date.isBefore(DateTime.now().subtract(Duration(minutes: 1)))) {
      int minutes = DateTime.now().difference(date).inMinutes;

      return '${minutes} min ago';
    }

    int seconds = DateTime.now().difference(date).inSeconds;
    return '${seconds} sec ago';
  }
}
