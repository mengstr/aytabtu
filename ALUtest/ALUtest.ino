#define MODE   2
#define CI     3
#define F2     4
#define A      5
#define F1     6
#define B      7
#define BINV   8
#define BCLR_  9
#define BSET_ 10
#define F2_   11
#define ENA   12
#define RESULT A0
#define COUT   A1


void setup() {
  Serial.begin(9600); 
  while (!Serial) {
    ; // wait for serial port to connect. Needed for Leonardo only
  }
  Serial.println("AYTABTU ALU Slice tester");
  
  pinMode(MODE,OUTPUT);
  pinMode(CI,OUTPUT);
  pinMode(F2,OUTPUT);
  pinMode(A,OUTPUT);
  pinMode(F1,OUTPUT);
  pinMode(B,OUTPUT);
  pinMode(BINV,OUTPUT);
  pinMode(BCLR_,OUTPUT);
  pinMode(BSET_,OUTPUT);
  pinMode(F2_,OUTPUT);
  pinMode(ENA,OUTPUT);
  pinMode(RESULT, INPUT);
  pinMode(COUT, INPUT);

  TestXor();
  TestOr();
  TestAnd();
  TestZero();
  TestOne();
  TestNot();
  TestInc();
  TestDec();
}


char *Scan(void) {
  static char table[1000];
  table[0]=0;

  strcat(table,"AB Ci- R Co\n");
  strcat(table,"===========\n");

  digitalWrite(A,LOW); digitalWrite(B,LOW); digitalWrite(CI,LOW);
  strcat(table,"00 0 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 
  
  digitalWrite(A,HIGH); digitalWrite(B,LOW); digitalWrite(CI,LOW);
  strcat(table,"10 0 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 

  digitalWrite(A,LOW); digitalWrite(B,HIGH); digitalWrite(CI,LOW);
  strcat(table,"01 0 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 

  digitalWrite(A,HIGH); digitalWrite(B,HIGH); digitalWrite(CI,LOW);
  strcat(table,"11 0 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 

  digitalWrite(A,LOW); digitalWrite(B,LOW); digitalWrite(CI,HIGH);
  strcat(table,"00 1 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 
  
  digitalWrite(A,HIGH); digitalWrite(B,LOW); digitalWrite(CI,HIGH);
  strcat(table,"10 1 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 

  digitalWrite(A,LOW); digitalWrite(B,HIGH); digitalWrite(CI,HIGH);
  strcat(table,"01 1 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 

  digitalWrite(A,HIGH); digitalWrite(B,HIGH); digitalWrite(CI,HIGH);
  strcat(table,"11 1 - ");
  if (digitalRead(RESULT)==HIGH) strcat(table,"1 "); else strcat(table,"0 "); 
  if (digitalRead(COUT)==HIGH) strcat(table,"1\n"); else strcat(table,"0\n"); 

  return table;
}


void TestXor(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  LOW);
  digitalWrite(F1,    HIGH);
  digitalWrite(F2,    HIGH);
  digitalWrite(F2_,   LOW);
  digitalWrite(BCLR_, HIGH);
  digitalWrite(BSET_, HIGH);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** XOR ***\n");
  Serial.println(pLine); 
}


void TestOr(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  LOW);
  digitalWrite(F1,    LOW);
  digitalWrite(F2,    HIGH);
  digitalWrite(F2_,   LOW);
  digitalWrite(BCLR_, HIGH);
  digitalWrite(BSET_, HIGH);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** OR ***\n");
  Serial.println(pLine); 
}

void TestAnd(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  LOW);
  digitalWrite(F1,    LOW);
  digitalWrite(F2,    LOW);
  digitalWrite(F2_,   HIGH);
  digitalWrite(BCLR_, HIGH);
  digitalWrite(BSET_, HIGH);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** AND ***\n");
  Serial.println(pLine); 
}


void TestZero(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  LOW);
  digitalWrite(F1,    LOW);
  digitalWrite(F2,    LOW);
  digitalWrite(F2_,   HIGH);
  digitalWrite(BCLR_, LOW);
  digitalWrite(BSET_, HIGH);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** ZERO ***\n");
  Serial.println(pLine); 
}

void TestOne(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  LOW);
  digitalWrite(F1,    LOW);
  digitalWrite(F2,    LOW);
  digitalWrite(F2_,   HIGH);
  digitalWrite(BCLR_, HIGH);
  digitalWrite(BSET_, LOW);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** ONE ***\n");
  Serial.println(pLine); 
}

void TestNot(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  LOW);
  digitalWrite(F1,    HIGH);
  digitalWrite(F2,    HIGH);
  digitalWrite(F2_,   LOW);
  digitalWrite(BCLR_, HIGH);
  digitalWrite(BSET_, LOW);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** ZERO ***\n");
  Serial.println(pLine); 
}


void TestInc(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  HIGH);
  digitalWrite(F1,    HIGH);
  digitalWrite(F2,    HIGH);
  digitalWrite(F2_,   LOW);
  digitalWrite(BCLR_, LOW);
  digitalWrite(BSET_, HIGH);
  digitalWrite(BINV,  HIGH);
  pLine=Scan();
  Serial.println("\n*** INC ***\n");
  Serial.println(pLine); 
}

void TestDec(void) {
  char *pLine;
  digitalWrite(ENA,   HIGH);
  digitalWrite(MODE,  HIGH);
  digitalWrite(F1,    HIGH);
  digitalWrite(F2,    HIGH);
  digitalWrite(F2_,   LOW);
  digitalWrite(BCLR_, HIGH);
  digitalWrite(BSET_, LOW);
  digitalWrite(BINV,  LOW);
  pLine=Scan();
  Serial.println("\n*** DEC ***\n");
  Serial.println(pLine); 
}


void loop() {
}
