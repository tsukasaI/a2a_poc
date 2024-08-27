import 'dart:async';
import 'dart:developer';
import 'package:flutter/material.dart';
import 'package:uni_links/uni_links.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';

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
        '/': (context) => const MyHomePage(title: 'A2A Demo Page'),
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
  String? name;
  String? loginChallenge;
  String? consentChallenge;

  @override
  void initState() {
    super.initState();
    initUniLinks();
  }

  Future<void> initUniLinks() async {
    linkStream.listen((String? link) {
      catchLink = link;
      getQueryParameters(link);
      setState(() {});
    }, onError: (err) {
      log(err);
    });
  }

  void getQueryParameters(String? link) {
    if (link == null) return;
    final uri = Uri.parse(link);
    name = uri.queryParameters['name'];
    loginChallenge = uri.queryParameters['login_challenge'];
    consentChallenge = uri.queryParameters['consent_challenge'];
  }

  Future<void> handleLogin(String locinCallengeCode) async {
    Uri uri =
        Uri(scheme: 'http', host: 'localhost', port: 3030, path: '/login');
    // if (!await launchUrl(uri)) {
    //   throw Exception('Could not launch $uri');
    // }

    http.Response res = await http.post(uri,
        body: <String, String>{'login_challenge': locinCallengeCode});
    String? location = res.headers['location'];
    String? cookie = res.headers['set-cookie'];
    if (location != null) {
      Uri redirectUri = Uri.parse(location);
      if (!await launchUrl(redirectUri)) {}
    }
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
            if (loginChallenge != null)
              ElevatedButton(
                  onPressed: () {
                    handleLogin(loginChallenge!);
                  },
                  child: const Text('Login')),
            ElevatedButton(
                onPressed: () {
                  // Navigator.pushNamed(context, '/second');
                },
                child: const Text('To Second')),
            const SizedBox(
              height: 20.0,
            ),
            Text(
              'name: $name',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            Text(
              'loginChallenge: $loginChallenge',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            Text(
              'consentChallenge: $consentChallenge',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ],
        ),
      ),
    );
  }
}
