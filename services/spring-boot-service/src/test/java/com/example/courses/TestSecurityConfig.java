package com.example.courses;

import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSSigner;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.gen.RSAKeyGenerator;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.jwt.JwtException;

import java.time.Instant;
import java.util.Collections;
import java.util.Map;

@Configuration
public class TestSecurityConfig {

    @Bean
    @Primary
    public JwtDecoder jwtDecoder() {
        return new JwtDecoder() {
            @Override
            public Jwt decode(String token) throws JwtException {
                try {
                    SignedJWT signedJWT = SignedJWT.parse(token);
                    JWTClaimsSet claims = signedJWT.getJWTClaimsSet();
                    JWSHeader header = signedJWT.getHeader();
                    
                    Map<String, Object> headers = new java.util.HashMap<>();
                    headers.put("alg", header.getAlgorithm().getName());
                    if (header.getKeyID() != null) {
                        headers.put("kid", header.getKeyID());
                    }
                    
                    return new Jwt(token, 
                            claims.getIssueTime() != null ? claims.getIssueTime().toInstant() : Instant.now(),
                            claims.getExpirationTime() != null ? claims.getExpirationTime().toInstant() : Instant.now().plusSeconds(3600),
                            headers,
                            claims.getClaims());
                } catch (Exception e) {
                    throw new JwtException("Failed to decode JWT: " + e.getMessage(), e);
                }
            }
        };
    }
}
