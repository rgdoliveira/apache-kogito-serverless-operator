package org.kie.kogito.postgresql.migrator;

import io.quarkus.test.Mock;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mockStatic;

public class DBConnectionCheckerTest {
    DBConnectionChecker dbConnectionChecker = new DBConnectionChecker();

    @Mock
    static DriverManager driverManager;

    @BeforeAll
    public static void init() {
        mockStatic(DriverManager.class);
    }

    @BeforeEach
    public void setupEach() {
        dbConnectionChecker.dataIndexDBURL = "jdbc:postgresql://db-service:5432/di";
        dbConnectionChecker.dataIndexDBUserName = "postgres";
        dbConnectionChecker.dataIndexDBPassword = "postgres";

        dbConnectionChecker.jobsServiceDBURL = "jdbc:postgresql://db-service:5432/js";
        dbConnectionChecker.jobsServiceDBUserName = "postgres";
        dbConnectionChecker.jobsServiceDBPassword = "postgres";
    }

    @Test
    public void testCheckDBConnections() throws SQLException {
        Mockito.when(driverManager.getConnection(anyString(), anyString(), anyString())).thenReturn(Mockito.mock(Connection.class));
        assertDoesNotThrow(() -> dbConnectionChecker.checkDataIndexDBConnection());
        assertDoesNotThrow(() -> dbConnectionChecker.checkJobsServiceDBConnection());
    }

    @Test
    public void testCheckDBConnectionsThrowSQLException() throws SQLException {
        Mockito.when(driverManager.getConnection(anyString(), anyString(), anyString())).thenThrow(SQLException.class);
        assertThrows(SQLException.class, () -> dbConnectionChecker.checkDataIndexDBConnection());
        assertThrows(SQLException.class, () -> dbConnectionChecker.checkJobsServiceDBConnection());
    }
}