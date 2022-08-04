import 'package:flutter/material.dart';
import 'package:yaru/yaru.dart';
import './theme.dart';
import './routes/setup.dart';
import 'package:yaru_icons/yaru_icons.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Builder(
        builder: (context) => YaruTheme(
          data: AppTheme.of(context),
          child: RouteSetup(title: "Setup"),
        ),
      ),
      debugShowCheckedModeBanner: false,
    );
  }
}
