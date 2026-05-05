package com.example.courses;

import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSSigner;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.gen.RSAKeyGenerator;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.Date;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
@Testcontainers
@Import(TestSecurityConfig.class)
class JwtAuthenticationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16")
            .withDatabaseName("courses")
            .withUsername("postgres")
            .withPassword("postgres");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
        // Configure Auth0 test settings
        registry.add("spring.security.oauth2.resourceserver.jwt.issuer-uri", () -> "https://test-auth0-tenant.auth0.com/");
    }

    @Autowired
    private MockMvc mockMvc;

    @Test
    void requestWithValidJwt_shouldReturn200() throws Exception {
        // Generate a test RSA key for JWT signing
        RSAKey rsaKey = new RSAKeyGenerator(2048).generate();
        JWSSigner signer = new RSASSASigner(rsaKey);

        JWTClaimsSet claims = new JWTClaimsSet.Builder()
                .subject("auth0|123456")
                .issuer("https://test-auth0-tenant.auth0.com/")
                .audience("test-audience")
                .expirationTime(new Date(System.currentTimeMillis() + 3600 * 1000))
                .claim("email", "test@example.com")
                .claim("name", "Test User")
                .build();

        SignedJWT signedJWT = new SignedJWT(new JWSHeader(com.nimbusds.jose.JWSAlgorithm.RS256), claims);
        signedJWT.sign(signer);

        String token = signedJWT.serialize();

        mockMvc.perform(get("/api/v1/courses")
                        .header("Authorization", "Bearer " + token))
                .andExpect(status().isOk());
    }
}
