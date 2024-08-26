import 'dart:async';
import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:uni_links/uni_links.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      routes: {
        '/': (context) => const MyHomePage(title: 'Flutter Demo Home Page'),
        '/first': (context) => const FirstScreen(),
        '/second': (context) => const SecondScreen(),
        '/test': (context) => const MyHomePage(title: 'Test deepling'),
      },
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  String? catchLink;
  String? parameter;

  @override
  void initState() {
    super.initState();
    initUniLinks();
  }

  Future<void> initUniLinks() async {
    linkStream.listen((String? link) {
      catchLink = link;
      parameter = getQueryParameter(link);
      setState(() {});
    }, onError: (err) {
      log(err);
    });
  }

  String? getQueryParameter(String? link) {
    if (link == null) return null;
    final uri = Uri.parse(link);
    String? name = uri.queryParameters['name'];
    return name;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Text(
              'catchLink: $catchLink',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            Text(
              'parameter: $parameter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            ElevatedButton(
                onPressed: () {
                  Navigator.pushNamed(context, '/first');
                },
                child: const Text('To First')),
            ElevatedButton(
                onPressed: () {
                  Navigator.pushNamed(context, '/second');
                },
                child: const Text('To Second')),
          ],
        ),
      ),
    );
  }
}
