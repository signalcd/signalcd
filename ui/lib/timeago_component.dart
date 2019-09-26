import 'dart:async';

import 'package:angular/angular.dart';

@Component(
  selector: 'signalcd-timeago',
  template: '<span>{{ ago }}</span>',
)
class TimeagoComponent implements OnInit, OnDestroy {
  Timer timer;
  String ago = '';

  @Input('date')
  DateTime date;

  @override
  void ngOnInit() {
    this.timer = Timer.periodic(const Duration(seconds: 1), tick);
    this.tick(this.timer);
  }

  @override
  void ngOnDestroy() {
    this.timer.cancel();
  }

  void tick(Timer timer) {
    this.ago = format(this.date);
  }

  String format(DateTime date) {
    if (date.isBefore(DateTime.now().subtract(Duration(days: 7)))) {
      return DatePipe().transform(date, 'mediumDate');
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
