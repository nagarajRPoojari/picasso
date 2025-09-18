ast.BlockStatement{
  Body: []ast.Statement{
    ast.ImportStatement{
      Name: "io",
      From: "io",
    },
    ast.VariableDeclarationStatement{
      Identifier: "PI",
      Constant: false,
      AssignedValue: ast.NumberExpression{
        Value: 3.14,
      },
      ExplicitType: ast.SymbolType{
        Value: "double",
      },
      IsStatic: false,
    },
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
            Value: "double",
          },
          IsStatic: false,
        },
        ast.VariableDeclarationStatement{
          Identifier: "x",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "Math",
              },
              Arguments: []ast.Expression{}, // p0
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "Math",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{}, // p1
          Name: "DirectoryReader",
          Body: []ast.Statement{
            ast.ExpressionStatement{
              Expression: ast.AssignmentExpression{
                Assignee: ast.MemberExpression{
                  Member: ast.SymbolExpression{
                    Value: "this",
                  },
                  Property: "y",
                },
                AssignedValue: ast.NumberExpression{
                  Value: 100.0,
                },
              },
            },
            ast.ExpressionStatement{
              Expression: ast.AssignmentExpression{
                Assignee: ast.MemberExpression{
                  Member: ast.MemberExpression{
                    Member: ast.SymbolExpression{
                      Value: "this",
                    },
                    Property: "x",
                  },
                  Property: "pi",
                },
                AssignedValue: ast.NumberExpression{
                  Value: 98.0,
                },
              },
            },
          },
          ReturnType: nil,
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: p1,
          Name: "math",
          Body: []ast.Statement{
            ast.ExpressionStatement{
              Expression: ast.CallExpression{
                Method: ast.MemberExpression{
                  Member: ast.SymbolExpression{
                    Value: "io",
                  },
                  Property: "printf",
                },
                Arguments: []ast.Expression{
                  ast.StringExpression{
                    Value: "this.y = %f  ",
                  },
                  ast.MemberExpression{
                    Member: ast.MemberExpression{
                      Member: ast.SymbolExpression{
                        Value: "this",
                      },
                      Property: "x",
                    },
                    Property: "pi",
                  },
                },
              },
            },
            ast.ExpressionStatement{
              Expression: ast.CallExpression{
                Method: ast.MemberExpression{
                  Member: ast.SymbolExpression{
                    Value: "this",
                  },
                  Property: "add",
                },
                Arguments: p0,
              },
            },
          },
          ReturnType: nil,
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: p1,
          Name: "add",
          Body: []ast.Statement{
            ast.ExpressionStatement{
              Expression: ast.CallExpression{
                Method: ast.MemberExpression{
                  Member: ast.SymbolExpression{
                    Value: "io",
                  },
                  Property: "printf",
                },
                Arguments: []ast.Expression{
                  ast.StringExpression{
                    Value: "this is inside math",
                  },
                },
              },
            },
          },
          ReturnType: nil,
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
            Value: 123.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "double",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: p1,
          Name: "Math",
          Body: []ast.Statement{}, // p2
          ReturnType: nil,
          IsStatic: false,
        },
      },
    },
    ast.FunctionDeclarationStatement{
      Parameters: p1,
      Name: "main",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "arr",
          Constant: false,
          AssignedValue: ast.ListExpression{
            Constants: []ast.Expression{
              ast.ListExpression{
                Constants: []ast.Expression{
                  ast.StringExpression{
                    Value: "1,",
                  },
                  ast.StringExpression{
                    Value: "2",
                  },
                  ast.StringExpression{
                    Value: "3",
                  },
                },
              },
              ast.ListExpression{
                Constants: []ast.Expression{
                  ast.StringExpression{
                    Value: "1,",
                  },
                  ast.StringExpression{
                    Value: "2",
                  },
                  ast.StringExpression{
                    Value: "3",
                  },
                },
              },
            },
          },
          ExplicitType: ast.ListType{
            Length: 10,
            Underlying: ast.ListType{
              Length: 2,
              Underlying: ast.SymbolType{
                Value: "string",
              },
            },
          },
          IsStatic: false,
        },
        ast.VariableDeclarationStatement{
          Identifier: "a",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "DirectoryReader",
              },
              Arguments: p0,
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "DirectoryReader",
          },
          IsStatic: false,
        },
        ast.VariableDeclarationStatement{
          Identifier: "n",
          Constant: false,
          AssignedValue: nil,
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.ExpressionStatement{
          Expression: ast.AssignmentExpression{
            Assignee: ast.SymbolExpression{
              Value: "n",
            },
            AssignedValue: ast.NumberExpression{
              Value: 278.0,
            },
          },
        },
        ast.ExpressionStatement{
          Expression: ast.CallExpression{
            Method: ast.MemberExpression{
              Member: ast.SymbolExpression{
                Value: "a",
              },
              Property: "math",
            },
            Arguments: p0,
          },
        },
        ast.VariableDeclarationStatement{
          Identifier: "z",
          Constant: false,
          AssignedValue: ast.StringExpression{
            Value: "hello world",
          },
          ExplicitType: ast.SymbolType{
            Value: "string",
          },
          IsStatic: false,
        },
        ast.IfStatement{
          Condition: ast.BinaryExpression{
            Left: ast.SymbolExpression{
              Value: "n",
            },
            Operator: lexer.Token{
              Kind: 19,
              Value: ">",
            },
            Right: ast.NumberExpression{
              Value: 10.0,
            },
          },
          Consequent: ast.BlockStatement{
            Body: []ast.Statement{
              ast.VariableDeclarationStatement{
                Identifier: "y",
                Constant: false,
                AssignedValue: ast.NumberExpression{
                  Value: 800.0,
                },
                ExplicitType: ast.SymbolType{
                  Value: "int",
                },
                IsStatic: false,
              },
              ast.IfStatement{
                Condition: ast.BinaryExpression{
                  Left: ast.NumberExpression{
                    Value: 100.0,
                  },
                  Operator: lexer.Token{
                    Kind: 19,
                    Value: ">",
                  },
                  Right: ast.NumberExpression{
                    Value: 20.0,
                  },
                },
                Consequent: ast.BlockStatement{
                  Body: []ast.Statement{
                    ast.ExpressionStatement{
                      Expression: ast.CallExpression{
                        Method: ast.MemberExpression{
                          Member: ast.SymbolExpression{
                            Value: "io",
                          },
                          Property: "printf",
                        },
                        Arguments: []ast.Expression{
                          ast.StringExpression{
                            Value: "value = %d ",
                          },
                          ast.SymbolExpression{
                            Value: "y",
                          },
                        },
                      },
                    },
                  },
                },
                Alternate: ast.BlockStatement{
                  Body: p2,
                },
              },
            },
          },
          Alternate: ast.BlockStatement{
            Body: []ast.Statement{
              ast.VariableDeclarationStatement{
                Identifier: "ni",
                Constant: false,
                AssignedValue: ast.NumberExpression{
                  Value: 200.0,
                },
                ExplicitType: ast.SymbolType{
                  Value: "int",
                },
                IsStatic: false,
              },
            },
          },
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
Duration: 771.75µs
PROCESSING === {{this} y} 
PROCESSING === {{{this} x} pi} 
PROCESSING === {{this} x} 
PROCESSING === {{{this} x} pi} 
PROCESSING === {{this} x} 
