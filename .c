ast.BlockStatement{
  Body: []ast.Statement{
    ast.ClassDeclarationStatement{
      Name: "DirectoryReader",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "x",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 200.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "float32",
          },
          IsStatic: true,
        },
        ast.VariableDeclarationStatement{
          Identifier: "y",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 1.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "float64",
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
          Name: "DirectoryReader",
          Body: []ast.Statement{
            ast.VariableDeclarationStatement{
              Identifier: "y",
              Constant: false,
              AssignedValue: ast.NumberExpression{
                Value: 100.0,
              },
              ExplicitType: nil,
              IsStatic: false,
            },
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.SymbolExpression{
                  Value: "y",
                },
              },
            },
          },
          ReturnType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: true,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "x",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
            ast.Parameter{
              Name: "y",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
          },
          Name: "sum",
          Body: []ast.Statement{
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.BinaryExpression{
                  Left: ast.SymbolExpression{
                    Value: "x",
                  },
                  Operator: lexer.Token{
                    Kind: 33,
                    Value: "+",
                  },
                  Right: ast.SymbolExpression{
                    Value: "y",
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
    ast.FunctionDeclarationStatement{
      Parameters: []ast.Parameter{},
      Name: "main",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "x",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "DirectoryReader",
              },
              Arguments: []ast.Expression{
                ast.NumberExpression{
                  Value: 200.0,
                },
              },
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "DirectoryReader",
          },
          IsStatic: false,
        },
        ast.ExpressionStatement{
          Expression: ast.CallExpression{
            Method: ast.MemberExpression{
              Member: ast.SymbolExpression{
                Value: "x",
              },
              Property: "sum",
            },
            Arguments: []ast.Expression{
              ast.NumberExpression{
                Value: 10.0,
              },
              ast.NumberExpression{
                Value: 10.0,
              },
            },
          },
        },
        ast.VariableDeclarationStatement{
          Identifier: "z",
          Constant: false,
          AssignedValue: ast.BinaryExpression{
            Left: ast.BinaryExpression{
              Left: ast.NumberExpression{
                Value: 10.0,
              },
              Operator: lexer.Token{
                Kind: 36,
                Value: "*",
              },
              Right: ast.NumberExpression{
                Value: 10.0,
              },
            },
            Operator: lexer.Token{
              Kind: 33,
              Value: "+",
            },
            Right: ast.NumberExpression{
              Value: 10.0,
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
      },
      ReturnType: nil,
      IsStatic: false,
    },
  },
}
Duration: 503.917µs
vars: map[]
classes: map[DirectoryReader:0x14000122cd0]
methods: map[DirectoryReader:0x14000122cd0]
