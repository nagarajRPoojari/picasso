ast.BlockStatement{
  Body: []ast.Statement{
    ast.ClassDeclarationStatement{
      Name: "DirectoryReader",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "y",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 112.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "x",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
          },
          Name: "sum",
          Body: []ast.Statement{
            ast.VariableDeclarationStatement{
              Identifier: "n",
              Constant: false,
              AssignedValue: ast.MemberExpression{
                Member: ast.SymbolExpression{
                  Value: "this",
                },
                Property: "y",
              },
              ExplicitType: ast.SymbolType{
                Value: "int",
              },
              IsStatic: false,
            },
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.BinaryExpression{
                  Left: ast.MemberExpression{
                    Member: ast.SymbolExpression{
                      Value: "this",
                    },
                    Property: "y",
                  },
                  Operator: lexer.Token{
                    Kind: 33,
                    Value: "+",
                  },
                  Right: ast.SymbolExpression{
                    Value: "x",
                  },
                },
              },
            },
          },
          ReturnType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
      },
    },
    ast.ClassDeclarationStatement{
      Name: "Math",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "pi",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 100.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
      },
    },
    ast.FunctionDeclarationStatement{
      Parameters: []ast.Parameter{},
      Name: "main",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "a",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "DirectoryReader",
              },
              Arguments: []ast.Expression{},
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "DirectoryReader",
          },
          IsStatic: false,
        },
        ast.VariableDeclarationStatement{
          Identifier: "z",
          Constant: false,
          AssignedValue: ast.CallExpression{
            Method: ast.MemberExpression{
              Member: ast.SymbolExpression{
                Value: "a",
              },
              Property: "sum",
            },
            Arguments: []ast.Expression{
              ast.NumberExpression{
                Value: 10.0,
              },
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.ReturnStatement{
          Value: ast.ExpressionStatement{
            Expression: ast.NumberExpression{
              Value: 0.0,
            },
          },
        },
      },
      ReturnType: ast.SymbolType{
        Value: "int32",
      },
      IsStatic: false,
    },
  },
}
Duration: 310.125µs
property -  y
assigning  &{i64 i64* %0 0}
property -  y
instance ---  &{DirectoryReader %DirectoryReader %DirectoryReader* %0}
assigning  &{DirectoryReader %DirectoryReader %DirectoryReader* %0}
assigning  &{i64 i64* %0 0}
instance - DirectoryReader
assigining - 112.000000 to  
 - a  
assigining - 0.000000 to  
 - n  
assigining - 0.000000 to  
 - z  
return z = 0.000000 1073741834
ran succesfully..
0
