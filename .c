ast.BlockStatement{
  Body: []ast.Statement{
    ast.ImportStatement{
      Name: "io",
      From: "builtin",
    },
    ast.ClassDeclarationStatement{
      Name: "String",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "s",
          Constant: false,
          AssignedValue: nil,
          ExplicitType: ast.SymbolType{
            Value: "string",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "s",
              Type: ast.SymbolType{
                Value: "string",
              },
            },
          },
          Name: "String",
          Body: []ast.Statement{
            ast.ExpressionStatement{
              Expression: ast.AssignmentExpression{
                Assignee: ast.MemberExpression{
                  Member: ast.SymbolExpression{
                    Value: "this",
                  },
                  Property: "s",
                },
                AssignedValue: ast.SymbolExpression{
                  Value: "s",
                },
              },
            },
          },
          ReturnType: nil,
          IsStatic: false,
        },
      },
    },
    ast.FunctionDeclarationStatement{
      Parameters: []ast.Parameter{},
      Name: "main",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "arr",
          Constant: false,
          AssignedValue: ast.ListExpression{
            Constants: []ast.Expression{
              ast.ListExpression{
                Constants: []ast.Expression{
                  ast.NewExpression{
                    Instantiation: ast.CallExpression{
                      Method: ast.SymbolExpression{
                        Value: "String",
                      },
                      Arguments: []ast.Expression{
                        ast.StringExpression{
                          Value: "hello",
                        },
                      },
                    },
                  },
                },
                EleType: &ast.ListType{ // p0
                  Length: 1,
                  Underlying: ast.SymbolType{
                    Value: "String",
                  },
                },
              },
              ast.ListExpression{
                Constants: []ast.Expression{
                  ast.NewExpression{
                    Instantiation: ast.CallExpression{
                      Method: ast.SymbolExpression{
                        Value: "String",
                      },
                      Arguments: []ast.Expression{
                        ast.StringExpression{
                          Value: "hello",
                        },
                      },
                    },
                  },
                },
                EleType: p0,
              },
              ast.ListExpression{
                Constants: []ast.Expression{
                  ast.NewExpression{
                    Instantiation: ast.CallExpression{
                      Method: ast.SymbolExpression{
                        Value: "String",
                      },
                      Arguments: []ast.Expression{
                        ast.StringExpression{
                          Value: "hello",
                        },
                      },
                    },
                  },
                },
                EleType: p0,
              },
              ast.ListExpression{
                Constants: []ast.Expression{
                  ast.NewExpression{
                    Instantiation: ast.CallExpression{
                      Method: ast.SymbolExpression{
                        Value: "String",
                      },
                      Arguments: []ast.Expression{
                        ast.StringExpression{
                          Value: "hello",
                        },
                      },
                    },
                  },
                },
                EleType: p0,
              },
            },
            EleType: &ast.ListType{
              Length: 4,
              Underlying: ast.ListType{
                Length: 1,
                Underlying: ast.SymbolType{
                  Value: "String",
                },
              },
            },
          },
          ExplicitType: ast.ListType{
            Length: 4,
            Underlying: ast.ListType{
              Length: 1,
              Underlying: ast.SymbolType{
                Value: "String",
              },
            },
          },
          IsStatic: false,
        },
        ast.ReturnStatement{
          Value: ast.ExpressionStatement{
            Expression: ast.NumberExpression{
              Value: 0.0,
            },
          },
          IsVoid: false,
        },
      },
      ReturnType: ast.SymbolType{
        Value: "int32",
      },
      IsStatic: false,
    },
  },
}
target --  String
target --  list.String[1]
