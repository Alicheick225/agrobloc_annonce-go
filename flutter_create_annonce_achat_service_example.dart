import 'dart:convert';
import 'package:http/http.dart' as http;

class AnnonceAchatService {
  final String baseUrl;

  AnnonceAchatService({required this.baseUrl});

  Future<void> createAnnonceAchat({
    required String userId,
    required String typeCultureId,
    required String statut,
    required String description,
    required double quantite,
    required double prixKg,
  }) async {
    final url = Uri.parse('\$baseUrl/annonces_achat');

    final Map<String, dynamic> body = {
      'user_id': userId,
      'type_culture_id': typeCultureId,
      'statut': statut,
      'description': description,
      'quantite': quantite,
      'prix_kg': prixKg,
    };

    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(body),
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to create annonce achat: \${response.body}');
    }
  }
}
