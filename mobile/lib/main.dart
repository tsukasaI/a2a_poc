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
  String? loginChallenge;
  String? consentChallenge;
  String? code;

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
      print(err);
    });
  }

  void getQueryParameters(String? link) {
    if (link == null) return;
    final uri = Uri.parse(link);
    loginChallenge = uri.queryParameters['login_challenge'];
    consentChallenge = uri.queryParameters['consent_challenge'];
    code = uri.queryParameters['code'];
  }

  Future<void> handleLogin(String locinCallengeCode) async {
    Uri uri =
        Uri(scheme: 'http', host: 'localhost', port: 3030, path: '/login');
    http.Response res = await http.post(uri, body: <String, String>{
      'login_challenge': locinCallengeCode,
      'username': 'testHydraUser0013',
    });
    String? consentUrl = res.headers['location'];
    if (consentUrl != null) {
      launchUrl(Uri.parse(consentUrl));
    }
  }

  Future<void> handleConsent(String consentCallengeCode) async {
    Uri uri =
        Uri(scheme: 'http', host: 'localhost', port: 3030, path: '/consent');
    http.Response res = await http.post(uri, body: <String, String>{
      'consent_challenge': consentCallengeCode,
      'consent': 'accept'
    });
    String? consentUrl = res.headers['location'];
    if (consentUrl != null) {
      launchUrl(Uri.parse(consentUrl), mode: LaunchMode.externalApplication);
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
            if (consentChallenge != null)
              ElevatedButton(
                  onPressed: () {
                    handleConsent(consentChallenge!);
                  },
                  child: const Text('Consent')),
            const SizedBox(
              height: 20.0,
            ),
            if (code != null)
              Text(
                'code: $code',
                style: Theme.of(context).textTheme.bodyMedium,
              ),
          ],
        ),
      ),
    );
  }
}
